package log

import (
	"os"
	"sync"

	srclog "github.com/sirupsen/logrus"
	config "github.com/zxz2801/gohttp_xxx/src/pkg/config"
)

// Log ...
type Log interface {
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
}

// glog ...
var glog = Log(srclog.StandardLogger())

// Logger ...
func Logger() Log {
	return glog
}

// InitLog ...
func InitLog() {
	llv := int(config.Global().Log.Level)
	lv := srclog.DebugLevel
	if len(srclog.AllLevels) > llv {
		lv = srclog.AllLevels[llv]
	}
	srclog.SetReportCaller(true)
	srclog.SetLevel(lv)
	srclog.SetFormatter(&srclog.TextFormatter{})
	srclog.SetOutput(os.Stdout)

	if config.Global().Log.Base != "" && config.Global().Log.File != "" {
		gwriter = &FileWriter{
			FilePath:    config.Global().Log.Base,
			FileName:    config.Global().Log.File,
			RotateSize:  uint64(config.Global().Log.RotateSize),
			RotateCount: int16(config.Global().Log.RotateCount),
			mux:         new(sync.Mutex),
		}

		err := gwriter.createNewFile()
		if err != nil {
			srclog.Errorf("createNewFile error[%s]", err.Error())
			return
		}
		srclog.SetOutput(gwriter)
	}

}
