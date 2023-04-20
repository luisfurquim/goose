package goose

import (
   "os"
   "fmt"
   "log"
   "reflect"
   "strings"
   "runtime"
   "log/syslog"
   "encoding/json"
)

// Alert is the basic type that implements the goose log object
type Alert uint8

type Geese map[string]interface{}

var GooseType  reflect.Type = reflect.PtrTo(reflect.TypeOf(Alert(0)))
var GooseValType  reflect.Type = reflect.TypeOf(Alert(0))

// SyslogGoose is a wrapper type around *syslog.Writer to ensure a syslogger that satisfies the io.Writer inrterface
type SyslogGoose struct {
   W *syslog.Writer
}

var notrace bool = true

// Write is the method that satisfies the io.Writer inrterface
func (ng SyslogGoose) Write(b []byte) (int, error) {
   return ng.W.Write(b)
}

func trace() string {
   var pc      []uintptr
   var path    []string
   var file      string
   var pkgfunc []string
   var pkg     []string
   var n         int
   var frames   *runtime.Frames
   var frame     runtime.Frame
   var id        string

   if notrace {
      return ""
   }
   pc = make([]uintptr, 10)  // at least 1 entry needed
   n = runtime.Callers(2, pc)
   if n == 0 {
      return ""
   }

   pc = pc[:n]
   frames = runtime.CallersFrames(pc)
   frame, _ = frames.Next()
   frame, _ = frames.Next()

   path = strings.Split(frame.File,string([]byte{os.PathSeparator}))
   file = path[len(path)-1]

   pkgfunc = strings.Split(frame.Function,string([]byte{os.PathSeparator}))
   pkgfunc = strings.Split(pkgfunc[len(pkgfunc)-1],".")
   pkg     = strings.Split(pkgfunc[0],"/")
   id      = fmt.Sprintf("{%s}[%s]<%s>(%d): ", pkg[len(pkg)-1], file, strings.Join(pkgfunc[1:],"."), frame.Line)
   return strings.Replace(strings.Replace(id, "%2e", ".", -1), "%", "%%", -1)
}

func deeptrace(stacklevel int) string {
   var pc      []uintptr
   var path    []string
   var file      string
   var pkgfunc []string
   var pkg     []string
   var n         int
   var frames   *runtime.Frames
   var frame     runtime.Frame
   var id        string

   if notrace {
      return ""
   }
   pc = make([]uintptr, 10)  // at least 1 entry needed
   n = runtime.Callers(2, pc)
   if n == 0 {
      return ""
   }

   pc = pc[:n]
   frames = runtime.CallersFrames(pc)
   frame, _ = frames.Next()
   frame, _ = frames.Next()
   for ; stacklevel > 0; stacklevel-- {
      frame, _ = frames.Next()
   }

   path = strings.Split(frame.File,string([]byte{os.PathSeparator}))
   file = path[len(path)-1]

   pkgfunc = strings.Split(frame.Function,string([]byte{os.PathSeparator}))
   pkgfunc = strings.Split(pkgfunc[len(pkgfunc)-1],".")
   pkg     = strings.Split(pkgfunc[0],"/")
   id      = fmt.Sprintf("{%s}[%s]<%s>(%d): ", pkg[len(pkg)-1], file, strings.Join(pkgfunc[1:],"."), frame.Line)
   return strings.Replace(strings.Replace(id, "%2e", ".", -1), "%", "%%", -1)
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

// TraceOn enables the inclusion of the package name, source file name, method/function caller and source line number in the log messages.
// As stated in https://golang.org/pkg/runtime/#Func.FileLine about the line numbering, "The result will not be accurate if pc is not a program
// counter within f".
func TraceOn() {
   notrace = false
}

// TraceOff disables the inclusion of the package name, source file name, method/function caller and source line number in the log messages.
// This is the default state of logging.
func TraceOff() {
   notrace = true
}

// Logf emits the messages based on the log.Printf
func (d Alert) Logf(level int, format string, parms ...interface{}) {
   if uint8(d) >= uint8(level) {
      log.Printf(trace() + format , parms...)
   }
}

// Fatalf behaves as Logf but stops program execution. The execution ends EVEN WHEN the log level of the message is higher the the current log level.
func (d Alert) Fatalf(level int, format string, parms ...interface{}) {
   if uint8(d) >= uint8(level) {
      log.Fatalf(trace() + format, parms...)
   }
   os.Exit(-1)
}

// Logf emits the messages based on fmt.Printf
func (d Alert) Printf(level int, format string, parms ...interface{}) {
   if uint8(d) >= uint8(level) {
      fmt.Printf(trace() + format, parms...)
   }
}

// Sprintf returns the messages as a string value
func (d Alert) Sprintf(level int, format string, parms ...interface{}) string {
   if uint8(d) >= uint8(level) {
      return fmt.Sprintf(trace() + format, parms...)
   }
   return ""
}

// Set accepts either integer or string values to initialize the goose Alert level
func (d *Alert) Set(level interface{}) {
   var n uint8

   switch level.(type) {
      case int8:
         (*d) = Alert(level.(int8))
      case int:
         (*d) = Alert(level.(int))
      case int16:
         (*d) = Alert(level.(int16))
      case int32:
         (*d) = Alert(level.(int32))
      case int64:
         (*d) = Alert(level.(int64))
      case uint8:
         (*d) = Alert(level.(uint8))
      case uint:
         (*d) = Alert(level.(uint))
      case uint16:
         (*d) = Alert(level.(uint16))
      case uint32:
         (*d) = Alert(level.(uint32))
      case uint64:
         (*d) = Alert(level.(uint64))
      case string:
         fmt.Sscanf(level.(string),"%d",&n)
         (*d) = Alert(n)
      case []byte:
         fmt.Sscanf(string(level.([]byte)),"%d",&n)
         (*d) = Alert(n)
   }
}

func (geese Geese) Set(level interface{}) {
   var i int
   var g interface{}
   var gtype reflect.Type
   var gval reflect.Value
   var valLevel uint64

	switch level.(type) {
	case int8, int16, int32, int64, int:
		valLevel = uint64(reflect.ValueOf(level).Int())
	case uint8, uint16, uint32, uint64, uint:
		valLevel = reflect.ValueOf(level).Uint()
	}

   for _, g = range geese {
      // We only consider parameters of struct pointer type
      // From these we set the goose.Alert fields

      // Ignore non pointer parameters
      if reflect.TypeOf(g).Kind() != reflect.Ptr {
         continue
      }

      // Ignore non struct pointer parameters
      gval = reflect.ValueOf(g).Elem()
      if gval.Kind() != reflect.Struct {
         continue
      }

      // Search for goose.Alert fields
      gtype = gval.Type()
      for i=0; i<gtype.NumField(); i++ {
         if gtype.Field(i).Type == GooseType.Elem() {
            gval.Field(i).SetUint(valLevel)
         }
      }
   }
}


func (geese Geese) UnmarshalJSON(data []byte) error {
   var m map[string]interface{}
   var err error
   var pkg string
   var rec interface{}
   var dstRec interface{}
   var gval reflect.Value
   var afldVal reflect.Value
   var recVal reflect.Value
   var ok bool
   var iter *reflect.MapIter

   log.Printf("---------------------\n")

   if err = json.Unmarshal(data, &m); err != nil {
      return err
   }

   for pkg, rec = range m["Goose"].(map[string]interface{}) {
      // Check if the key in the source exists in destiny
      if dstRec, ok = geese[pkg]; !ok {
         continue
      }

      // We only consider parameters of struct pointer type
      // From these we set the goose.Alert fields

      // Ignore non pointer parameters
      if reflect.TypeOf(dstRec).Kind() != reflect.Ptr {
         continue
      }

      // Ignore non struct pointer parameters
      gval = reflect.ValueOf(dstRec).Elem()
      if gval.Kind() != reflect.Struct {
         continue
      }

      // Search for goose.Alert fields
      recVal = reflect.ValueOf(rec)
      iter = recVal.MapRange()
      for iter.Next() {
         afldVal = gval.FieldByName(iter.Key().String())
         if !afldVal.IsZero() && afldVal.Type() == GooseType.Elem() {
            afldVal.Set(iter.Value().Elem().Convert(GooseValType))
         }
      }
   }

   return nil
}


func (geese Geese) Get() map[string]interface{} {
   return map[string]interface{}(geese)
}

// Logf emits the messages based on the log.Printf
func (d Alert) DeepLogf(stacklevel, level int, format string, parms ...interface{}) {
   if uint8(d) >= uint8(level) {
      log.Printf(deeptrace(stacklevel) + format , parms...)
   }
}

// Fatalf behaves as Logf but stops program execution. The execution ends EVEN WHEN the log level of the message is higher the the current log level.
func (d Alert) DeepFatalf(stacklevel, level int, format string, parms ...interface{}) {
   if uint8(d) >= uint8(level) {
      log.Fatalf(deeptrace(stacklevel) + format, parms...)
   }
   os.Exit(-1)
}

// Logf emits the messages based on fmt.Printf
func (d Alert) DeepPrintf(stacklevel, level int, format string, parms ...interface{}) {
   if uint8(d) >= uint8(level) {
      fmt.Printf(deeptrace(stacklevel) + format, parms...)
   }
}

// Sprintf returns the messages as a string value
func (d Alert) DeepSprintf(stacklevel, level int, format string, parms ...interface{}) string {
   if uint8(d) >= uint8(level) {
      return fmt.Sprintf(deeptrace(stacklevel) + format, parms...)
   }
   return ""
}
