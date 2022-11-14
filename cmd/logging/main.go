package main

import (
	"github.com/030/logging/pkg/logging"
	log "github.com/sirupsen/logrus"
)

func main() {
	l := logging.Logging{File: "blah.log", Level: "trace", Syslog: true}
	f, err := l.Setup()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	log.Trace("trace")
	log.Debug("debug")
	log.Info("info")
	log.Warn("warn")
	log.Error("error")
}
