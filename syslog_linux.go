package goose

import (
   "os"
   "log"
   "strings"
   "log/syslog"
)

// SyslogGoose is a wrapper type around *syslog.Writer to ensure a syslogger that satisfies the io.Writer inrterface
type SyslogGoose struct {
   W *syslog.Writer
}

// Write is the method that satisfies the io.Writer inrterface
func (ng SyslogGoose) Write(b []byte) (int, error) {
   return ng.W.Write(b)
}

// UseSyslogNet redirects the log output from os.Stderr to the system logger
// connecting to it via network.
func UseSyslogNet(proto string, srv string, priority syslog.Priority) error {
   var logOutput     SyslogGoose
   var binParts    []string
   var binName       string
   var err           error

   binParts = strings.Split(os.Args[0],string([]byte{os.PathSeparator}))
   binName  = binParts[len(binParts)-1]

   logOutput.W, err = syslog.Dial(proto, srv, priority, binName)
   if err != nil {
      return err
   }
   log.SetOutput(logOutput)
   log.SetFlags(0)
   return nil
}

