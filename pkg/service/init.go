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

//const (
//	connHost = "192.168.100.83"
//	connPort = "8080"
//	connType = "tcp"
//)

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

		// Respond to client
		response := []byte("Message received")
		conn.Write(response)
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

	//Test OctetString String
	receivedBytes := []byte{126, 0, 86, 1, 2, 0, 0, 33, 88, 232, 229, 214, 85, 45, 81, 244, 221, 211, 121, 161, 48, 159, 98, 68, 32, 16, 209, 220, 248, 180, 89, 194, 128, 0, 187, 253, 229, 212, 166, 22, 70, 23, 141, 116, 224, 208, 225, 177, 128, 0, 185, 53, 48, 220, 24, 13, 170, 65, 17, 68, 212, 21, 238, 32, 114, 63, 212, 30, 87, 226, 60, 141, 85, 63, 0, 96, 103, 219, 67, 195, 128, 0, 68, 41, 240, 208, 251, 138, 95, 233, 43, 138, 188, 195, 64, 166, 234, 195, 155, 182, 178, 5, 64, 100, 0, 211, 207, 131, 139, 144, 147, 131, 128, 0, 161, 43, 126, 74, 48, 115, 228, 240, 51, 216, 80, 224, 20, 204, 203, 224, 232, 40, 166, 25, 108, 84, 252, 47, 157, 223, 46, 1, 24, 44, 128, 0, 244, 210, 186, 162, 42, 102, 4, 168, 215, 139, 87, 185, 5, 146, 166, 232, 120, 118, 190, 156, 97, 156, 188, 36, 88, 72, 162, 8, 106, 6, 128, 0, 64, 67, 55, 173, 233, 224, 165, 96, 133, 115, 12, 127, 165, 123, 161, 194, 195, 138, 94, 153, 217, 162, 182, 137, 71, 32, 241, 51, 182, 52, 128, 0, 146, 37, 115, 217, 179, 93, 178, 46, 38, 24, 143, 112, 206, 178, 172, 252, 94, 98, 152, 126, 243, 65, 104, 80, 228, 108, 40, 255, 140, 192, 128, 0, 21, 119, 195, 215, 148, 31, 159, 72, 73, 90, 163, 253, 7, 195, 88, 92, 7, 101, 232, 199, 237, 165, 195, 197, 20, 164, 218, 12, 223, 138, 128, 0, 94, 27, 146, 170, 60, 244, 95, 85, 95, 80, 245, 74, 16, 106, 94, 113, 143, 80, 242, 228, 41, 102, 33, 154, 225, 187, 102, 156, 97, 72, 128, 0, 101, 111, 193, 226, 194, 140, 89, 152, 128, 170, 216, 209, 229, 188, 81, 94, 168, 128, 173, 163, 221, 221, 205, 107, 212, 98, 124, 223, 179, 237, 128, 0, 17, 217, 20, 126, 75, 15, 139, 4, 23, 14, 55, 156, 25, 252, 195, 230, 112, 156, 99, 214, 152, 102, 169, 175}

	//Test OctetString Ending

	// Listen for incoming connections on port 8080
	ln, err := net.Listen("tcp", ":12345")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer ln.Close()

	fmt.Println("Server is listening on port 12345")

	//go func() {
	//	for {
	//		// Accept incoming connection
	//		conn, err := ln.Accept()
	//		if err != nil {
	//			fmt.Println("Error accepting connection:", err.Error())
	//			return
	//		}
	//
	//		fmt.Println("New client connected:", conn.RemoteAddr())
	//
	//		// Handle connections concurrently
	//		go handleConnection(conn)
	//	}
	//}()

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
