// Package goose is a simple golang log package for debug purposes.
// The idea is very simple. You create variables of type goose.Alert and set its debug level.
// We suggest you set these variables to level 0 when you don't want any log messages at all.
// Set to level 1 to log error messages. Set to level 2 or above for increasing levels of log
// verbosity.

The package offers 4 methods to emit debug messages:

```Go
func (d Alert) Logf(level int, format string, parms ...interface{})
```

Based on the log.Printf.


```Go
func (d Alert) Fatalf(level int, format string, parms ...interface{})
```

Same as above but stops program execution. The execution ends EVEN WHEN the log level of the message is higher the the current log level.


```Go
func (d Alert) Printf(level int, format string, parms ...interface{})
```

Based on fmt.Printf.


```Go
func (d Alert) Sprintf(level int, format string, parms ...interface{}) string
```

Based on fmt.Sprintf.




In all the above methods the level parameter determines the log level of the message. The message will be actually emited only if, when the method is called, the log level is equal or higher than the message's log level. An empty string is returned by the Sprintf method if called with a message log level higher than the current log level.


## Example:

In your initialization code, you may set the log levels of debug messages:

```Go
   ...

   reptilian.Goose  = goose.Alert(1) // Only error messages will be logged

   reptext.Goose    = goose.Alert(0) // No debug messages at all

   djparser.Goose   = goose.Alert(4) // Error messages and less verbose messages (levels 2~4) will be logged

   ...
```

In the package to be debugged:

```Go

import (
   ...

   "github.com/luisfurquim/goose"
   ...
)

...

var Goose goose.Alert // Exported symbol needed only if you want to allow external control of the debug level

...

   // throughout your code you may log anything you want at varying levels of importance

   Goose.Logf(3, "Final Off=%d (%o)", d.Buf.Off, d.Buf.Off)

   ...

   Goose.Logf(7, "Index: %#v", d.Index)

   // Logs will be actually printed only if the first parameter (the message's log level) is lower or equal than the current log level indicated by the Goose variable. Remember to never use the zero value, like Goose.Logf(0,...), as we want to make the log level 0 to print no debug messages at all.

```




You may set multiple loggers if you need finer control on the verbosity

```Go
type T1 struct {

...

}


type T2 struct {

...

}

var GooseT1 goose.Alert
var GooseT2 goose.Alert

   ...

   GooseT1 = goose.Alert(1) // Only error messages will be logged
   GooseT2 = goose.Alert(3) // More verbosity...

   ...

   GooseT1.Logf(2, "Final Off=%d (%o)", d.Buf.Off, d.Buf.Off) // not printed

   ...

   GooseT2.Logf(2, "Index: %#v", d.Index) // printed

```

You may add/remove source code reference by enabling/disabling trace.
With trace enabled Goose automatically adds source code reference in
the format {package}[source filename]&lt;function/method&gt;(line number).

Check it with the following code:


```Go

   Goose = goose.Alert(1)
   Goose.Logf(1,"no trace")
   goose.TraceOn()
   Goose.Logf(1,"trace")
   goose.TraceOff()
   Goose.Logf(1,"no trace")

```

You may redirect the output to the syslogger just calling UseSyslogNet:


```Go

   goose.UseSyslogNet("tcp", "myloghost.mydomain:514", syslog.LOG_ERR|syslog.LOG_LOCAL7)

```

