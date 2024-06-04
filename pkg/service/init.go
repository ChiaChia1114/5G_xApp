package service

import (
	"fmt"
	"github.com/free5gc/util/httpwrapper"
	logger_util "github.com/free5gc/util/logger"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"
	xApp_context "xApp/internal/context"
	"xApp/internal/logger"
	"xApp/internal/util"
	"xApp/pkg/factory"
	context "xApp/pkg/service/context"
	Authtimer "xApp/pkg/service/timer"
)

type XApp struct {
	KeyLogPath string
}

type (
	Commands struct {
		config string
	}
)

var commands Commands

var cliCmd = []cli.Flag{
	cli.StringFlag{
		Name:  "config, c",
		Usage: "Load configuration from `FILE`",
	},
	cli.StringFlag{
		Name:  "log, l",
		Usage: "Output NF log to `FILE`",
	},
	cli.StringFlag{
		Name:  "log5gc, lc",
		Usage: "Output xApp log to `FILE`",
	},
}

func (*XApp) GetCliCmd() (flags []cli.Flag) {
	return cliCmd
}

func (xApp *XApp) Initialize(c *cli.Context) error {
	commands = Commands{
		config: c.String("config"),
	}

	if commands.config != "" {
		if err := factory.InitConfigFactory(commands.config); err != nil {
			return err
		}
	} else {
		if err := factory.InitConfigFactory(util.XAppDefaultConfigPath); err != nil {
			return err
		}
	}

	if err := factory.CheckConfigVersion(); err != nil {
		return err
	}

	if _, err := factory.XAppConfig.Validate(); err != nil {
		return err
	}

	xApp.setLogLevel()

	return nil
}

func (xApp *XApp) setLogLevel() {
	if factory.XAppConfig.Logger == nil {
		logger.InitLog.Warnln("xApp config without log level setting!!!")
		return
	}

	if factory.XAppConfig.Logger.XAPP != nil {
		if factory.XAppConfig.Logger.XAPP.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.XAppConfig.Logger.XAPP.DebugLevel); err != nil {
				logger.InitLog.Warnf("xApp Log level [%s] is invalid, set to [info] level",
					factory.XAppConfig.Logger.XAPP.DebugLevel)
				logger.SetLogLevel(logrus.InfoLevel)
			} else {
				logger.InitLog.Infof("xApp Log level is set to [%s] level", level)
				logger.SetLogLevel(level)
			}
		} else {
			logger.InitLog.Warnln("xApp Log level not set. Default set to [info] level")
			logger.SetLogLevel(logrus.InfoLevel)
		}
		logger.SetReportCaller(factory.XAppConfig.Logger.XAPP.ReportCaller)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Buffer to hold incoming data
	buffer := make([]byte, 1024)

	for {
		// Read incoming data into buffer
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return
		}
		octetString := buffer[:bytesRead]
		//fmt.Println("Received OctetString:", octetString)
		if octetString != nil {

			OriginalNASMessage := HandleMessageSelection(octetString)
			// Respond to client
			//fmt.Println("OriginalNASMessage: ", OriginalNASMessage)
			_, err = conn.Write(OriginalNASMessage)
			if err != nil {
				fmt.Println("Error writing response:", err.Error())
			}

			//HandleOtherMessage
			UEid := 1
			UE_status, result := context.GetCountByUEid(UEid)
			if result != true {
				fmt.Println("Failed to get UE count")
			} else {
				if UE_status == 0 {
					HandleAuthenticationVectorsPreparetion()
				}
			}

		}
		//_, err = conn.Write("No data")
		//if err != nil {
		//	fmt.Println("Error writing response:", err.Error())
		//}

	}
}

func (xApp *XApp) Start() {
	logger.InitLog.Infoln("Server started")
	router := logger_util.NewGinWithLogrus(logger.GinLog)

	xApp_context.Init()
	self := xApp_context.GetSelf()
	addr := fmt.Sprintf("%s:%d", self.BindingIPv4, self.SBIPort)

	// Listen for incoming connections on port 8080
	ln, err := net.Listen("tcp", ":12345")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer ln.Close()
	fmt.Println("Server is listening on port 12345")

	_, err = context.GenerateToken()
	if err != nil {
		fmt.Println("Error GenerateToken:", err.Error())
	}

	// Terry Modify start: Add Timer to calculate service time
	StartTime := time.Now()
	TimernewUe := Authtimer.NewServiceTimer(1, StartTime)
	Authtimer.StoreTimeStamp(TimernewUe)
	// Terry Modify end: Add Timer to calculate service time

	go func() {
		for {
			// Accept incoming connection
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("Error accepting connection:", err.Error())
				return
			}

			fmt.Println("New client connected:", conn.RemoteAddr())

			// Handle connections concurrently
			go handleConnection(conn)
		}
	}()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				// Print stack for panic to log. Fatalf() will let program exit.
				logger.InitLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
			}
		}()

		<-signalChannel
		xApp.Terminate()
		os.Exit(0)
	}()

	server, err := httpwrapper.NewHttp2Server(addr, xApp.KeyLogPath, router)
	if server == nil {
		logger.InitLog.Errorf("Initialize HTTP server failed: %+v", err)
		return
	}

	if err != nil {
		logger.InitLog.Warnf("Initialize HTTP server: +%v", err)
	}

	serverScheme := factory.XAppConfig.Configuration.Sbi.Scheme
	if serverScheme == "http" {
		err = server.ListenAndServe()
	}

	if err != nil {
		logger.InitLog.Fatalf("HTTP server setup failed: %+v", err)
	}
}

func (xApp *XApp) Terminate() {
	logger.InitLog.Infof("Terminating AUSF...")
	logger.InitLog.Infof("AUSF terminated")
}
