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
	xApp_context "xApp/internal/context"
	"xApp/internal/logger"
	"xApp/internal/util"
	"xApp/pkg/factory"
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
		fmt.Println("Received OctetString:", octetString)
		OriginalNASMessage, OtherNASMessage := HandleOctetString(octetString)
		// Respond to client
		response := string(OriginalNASMessage)
		//conn.Write(response)
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing response:", err.Error())
		}

		//HandleOtherMessage
		HandleOtherMessage(OtherNASMessage)
	}
}

func (xApp *XApp) Start() {
	logger.InitLog.Infoln("Server started")
	router := logger_util.NewGinWithLogrus(logger.GinLog)

	// Setting TLS initailize
	//ueauthentication.AddService(router)
	//
	//pemPath := util.AusfDefaultPemPath
	//keyPath := util.AusfDefaultKeyPath
	//sbi := factory.AusfConfig.Configuration.Sbi
	//if sbi.Tls != nil {
	//	pemPath = sbi.Tls.Pem
	//	keyPath = sbi.Tls.Key
	//}
	xApp_context.Init()
	fmt.Println("Hello world! 4")
	self := xApp_context.GetSelf()
	// Register to NRF
	//profile, err := consumer.BuildNFInstance(self)
	//if err != nil {
	//	logger.InitLog.Error("Build AUSF Profile Error")
	//}
	//_, self.NfId, err = consumer.SendRegisterNFInstance(self.NrfUri, self.NfId, profile)
	//if err != nil {
	//	logger.InitLog.Errorf("AUSF register to NRF Error[%s]", err.Error())
	//}

	addr := fmt.Sprintf("%s:%d", self.BindingIPv4, self.SBIPort)

	// Listen for incoming connections on port 8080
	ln, err := net.Listen("tcp", ":12345")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer ln.Close()

	fmt.Println("Server is listening on port 12345")

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
	//} else if serverScheme == "https" {
	//	err = server.ListenAndServeTLS(pemPath, keyPath)
	//}

	if err != nil {
		logger.InitLog.Fatalf("HTTP server setup failed: %+v", err)
	}
}

func (xApp *XApp) Terminate() {
	logger.InitLog.Infof("Terminating AUSF...")
	// deregister with NRF
	//problemDetails, err := consumer.SendDeregisterNFInstance()
	//if problemDetails != nil {
	//	logger.InitLog.Errorf("Deregister NF instance Failed Problem[%+v]", problemDetails)
	//} else if err != nil {
	//	logger.InitLog.Errorf("Deregister NF instance Error[%+v]", err)
	//} else {
	//	logger.InitLog.Infof("Deregister from NRF successfully")
	//}
	logger.InitLog.Infof("AUSF terminated")
}
