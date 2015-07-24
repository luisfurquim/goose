package goose

import (
   "os"
   "fmt"
   "log"
   "strings"
   "runtime"
   "log/syslog"
)

type Alert uint8


type SyslogGoose struct {
   W *syslog.Writer
}

var notrace bool = true

func (ng SyslogGoose) Write(b []byte) (int, error) {
   return ng.W.Write(b)
}

func trace() string {
   var pc      []uintptr
   var f        *runtime.Func
   var path    []string
   var file      string
   var line      int
   var pkgfunc []string
   var pkg     []string

   if notrace {
      return ""
   }
   pc = make([]uintptr, 10)  // at least 1 entry needed
   runtime.Callers(2, pc)
   f = runtime.FuncForPC(pc[1])
   file, line = f.FileLine(pc[1])
   path = strings.Split(file,string([]byte{os.PathSeparator}))
   file = path[len(path)-1]
   pkgfunc = strings.Split(f.Name(),string([]byte{os.PathSeparator}))
   pkgfunc = strings.Split(pkgfunc[len(pkgfunc)-1],".")
   pkg     = strings.Split(pkgfunc[0],"/")
   return fmt.Sprintf("{%s}[%s]<%s>(%d): ", pkg[len(pkg)-1], file, strings.Join(pkgfunc[1:],"."), line)
}

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

func TraceOn() {
   notrace = false
}


func TraceOff() {
   notrace = true
}


func (d Alert) Logf(level int, format string, parms ...interface{}) {
   if uint8(d) >= uint8(level) {
      log.Printf(trace() + format , parms...)
   }
}

func (d Alert) Fatalf(level int, format string, parms ...interface{}) {
   if uint8(d) >= uint8(level) {
      log.Fatalf(trace() + format, parms...)
   }
   os.Exit(-1)
}

func (d Alert) Printf(level int, format string, parms ...interface{}) {
   if uint8(d) >= uint8(level) {
      fmt.Printf(trace() + format, parms...)
   }
}

func (d Alert) Sprintf(level int, format string, parms ...interface{}) string {
   if uint8(d) >= uint8(level) {
      return fmt.Sprintf(trace() + format, parms...)
   }
   return ""
}
