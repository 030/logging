package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	jww "github.com/spf13/jwalterweatherman"
)

type Logging struct {
	File, Level string
	Syslog      bool
}

func file(path string) (*os.File, error) {
	f, err := os.OpenFile(filepath.Clean(path), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return nil, err
	}

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
	})

	log.AddHook(&writer.Hook{
		Writer:    f,
		LogLevels: []log.Level{log.TraceLevel, log.DebugLevel, log.InfoLevel, log.WarnLevel, log.ErrorLevel},
	})

	return f, err
}

func fileOrSyslog(path string, syslog bool) (*os.File, error) {
	var f *os.File
	var err error
	if path != "" {
		f, err = file(path)
		if err != nil {
			return nil, err
		}
	}

	if syslog {
		if runtime.GOOS != "windows" {
			syslogFile()
		}
	}

	return f, nil
}

func (l *Logging) Setup() (*os.File, error) {
	log.SetReportCaller(true)

	switch l.Level {
	case "trace":
		jww.SetLogThreshold(jww.LevelTrace)
		jww.SetStdoutThreshold(jww.LevelTrace)
		log.SetLevel(log.TraceLevel)
	case "debug":
		jww.SetLogThreshold(jww.LevelDebug)
		jww.SetStdoutThreshold(jww.LevelDebug)
		log.SetLevel(log.DebugLevel)
	case "info":
		jww.SetLogThreshold(jww.LevelInfo)
		jww.SetStdoutThreshold(jww.LevelInfo)
		log.SetLevel(log.InfoLevel)
	case "warn":
		jww.SetLogThreshold(jww.LevelWarn)
		jww.SetStdoutThreshold(jww.LevelWarn)
		log.SetLevel(log.WarnLevel)
	case "error":
		jww.SetLogThreshold(jww.LevelError)
		jww.SetStdoutThreshold(jww.LevelError)
		log.SetLevel(log.ErrorLevel)
	case "none":
		log.SetOutput(io.Discard)
	default:
		return nil, fmt.Errorf("logLevel: '%s' is invalid and should have been 'trace',"+
			" 'debug', 'info', 'warn', 'error' or 'none'", l.Level)
	}

	f, err := fileOrSyslog(l.File, l.Syslog)
	if err != nil {
		return nil, err
	}

	return f, nil
}
