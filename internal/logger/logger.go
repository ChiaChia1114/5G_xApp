package logger

import (
	"github.com/asaskevich/govalidator"
	"reflect"
	"time"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

var (
	log                 *logrus.Logger
	AppLog              *logrus.Entry
	InitLog             *logrus.Entry
	CfgLog              *logrus.Entry
	UeAuthPostLog       *logrus.Entry
	Auth5gAkaComfirmLog *logrus.Entry
	EapAuthComfirmLog   *logrus.Entry
	HandlerLog          *logrus.Entry
	ContextLog          *logrus.Entry
	ConsumerLog         *logrus.Entry
	GinLog              *logrus.Entry
)

type Logger struct {
	XAPP *LogSetting `yaml:"XAPP" valid:"optional"`
}

func init() {
	log = logrus.New()
	log.SetReportCaller(false)

	log.Formatter = &formatter.Formatter{
		TimestampFormat: time.RFC3339,
		TrimMessages:    true,
		NoFieldsSpace:   true,
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
	}

	AppLog = log.WithFields(logrus.Fields{"component": "XAPP", "category": "App"})
	InitLog = log.WithFields(logrus.Fields{"component": "XAPP", "category": "Init"})
	CfgLog = log.WithFields(logrus.Fields{"component": "XAPP", "category": "CFG"})
	UeAuthPostLog = log.WithFields(logrus.Fields{"component": "XAPP", "category": "UeAuthPost"})
	Auth5gAkaComfirmLog = log.WithFields(logrus.Fields{"component": "XAPP", "category": "5gAkaAuth"})
	EapAuthComfirmLog = log.WithFields(logrus.Fields{"component": "XAPP", "category": "EapAkaAuth"})
	HandlerLog = log.WithFields(logrus.Fields{"component": "XAPP", "category": "Handler"})
	ContextLog = log.WithFields(logrus.Fields{"component": "XAPP", "category": "ctx"})
	ConsumerLog = log.WithFields(logrus.Fields{"component": "XAPP", "category": "Consumer"})
	GinLog = log.WithFields(logrus.Fields{"component": "XAPP", "category": "GIN"})
}

//func LogFileHook(logNfPath string, log5gcPath string) error {
//	if fullPath, err := logger_util.CreateFree5gcLogFile(log5gcPath); err == nil {
//		if fullPath != "" {
//			free5gcLogHook, hookErr := logger_util.NewFileHook(fullPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
//			if hookErr != nil {
//				return hookErr
//			}
//			log.Hooks.Add(free5gcLogHook)
//		}
//	} else {
//		return err
//	}
//
//	if fullPath, err := logger_util.CreateNfLogFile(logNfPath, "ausf.log"); err == nil {
//		selfLogHook, hookErr := logger_util.NewFileHook(fullPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
//		if hookErr != nil {
//			return hookErr
//		}
//		log.Hooks.Add(selfLogHook)
//	} else {
//		return err
//	}
//
//	return nil
//}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(enable bool) {
	log.SetReportCaller(enable)
}

func (l *Logger) Validate() (bool, error) {
	logger := reflect.ValueOf(l).Elem()
	for i := 0; i < logger.NumField(); i++ {
		if logSetting := logger.Field(i).Interface().(*LogSetting); logSetting != nil {
			result, err := logSetting.validate()
			return result, err
		}
	}

	result, err := govalidator.ValidateStruct(l)
	return result, err
}

type LogSetting struct {
	DebugLevel   string `yaml:"debugLevel" valid:"debugLevel"`
	ReportCaller bool   `yaml:"ReportCaller" valid:"type(bool)"`
}

func (l *LogSetting) validate() (bool, error) {
	govalidator.TagMap["debugLevel"] = govalidator.Validator(func(str string) bool {
		if str == "panic" || str == "fatal" || str == "error" || str == "warn" ||
			str == "info" || str == "debug" || str == "trace" {
			return true
		} else {
			return false
		}
	})

	result, err := govalidator.ValidateStruct(l)
	return result, err
}
