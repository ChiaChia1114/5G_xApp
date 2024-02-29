package main

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/free5gc/util/version"
	"github.com/urfave/cli"
	"os"
	"xApp/internal/logger"
	"xApp/pkg/service"
)

var xApp = &service.XApp{}

func main() {

	app := cli.NewApp()
	app.Name = "xApp"
	app.Usage = "RT Lab xApp-based AKA test"
	app.Action = action
	app.Flags = xApp.GetCliCmd()
	if err := app.Run(os.Args); err != nil {
		logger.AppLog.Errorf("xApp Run error: %v\n", err)
		return
	}

	//m := nas.NewMessage()
	//m.GmmMessage = nas.NewGmmMessage()
}

func action(c *cli.Context) error {
	//if err := initLogFile(c.String("log"), c.String("log5gc")); err != nil {
	//	logger.AppLog.Errorf("%+v", err)
	//	return err
	//}

	if err := xApp.Initialize(c); err != nil {
		switch errType := err.(type) {
		case govalidator.Errors:
			validErrs := err.(govalidator.Errors).Errors()
			for _, validErr := range validErrs {
				logger.CfgLog.Errorf("%+v", validErr)
			}
		default:
			logger.CfgLog.Errorf("%+v", errType)
		}
		logger.CfgLog.Errorf("[-- PLEASE REFER TO SAMPLE CONFIG FILE COMMENTS --]")
		return fmt.Errorf("Failed to initialize !!")
	}

	logger.AppLog.Infoln(c.App.Name)
	logger.AppLog.Infoln("xApp version: ", version.GetVersion())

	xApp.Start()

	return nil
}

//func initLogFile(logNfPath, log5gcPath string) error {
//	xApp.KeyLogPath = util.XAppDefaultKeyLogPath
//
//	//if err := logger.LogFileHook(logNfPath, log5gcPath); err != nil {
//	//	return err
//	//}
//
//	if logNfPath != "" {
//		nfDir, _ := filepath.Split(logNfPath)
//		tmpDir := filepath.Join(nfDir, "key")
//		if err := os.MkdirAll(tmpDir, 0775); err != nil {
//			logger.InitLog.Errorf("Make directory %s failed: %+v", tmpDir, err)
//			return err
//		}
//		_, name := filepath.Split(util.XAppDefaultKeyLogPath)
//		xApp.KeyLogPath = filepath.Join(tmpDir, name)
//	}
//
//	return nil
//}
