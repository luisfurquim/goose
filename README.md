# goose
Simple golang log package for debug purposes

The idea is very simple. You create variables of type goose.Alert and set its debug level.
We suggest you set these variables to level 0 when you don't want any log messages at all. 
Set to level 1 to log error messages. Set to level 2 or above for increasing levels of log
verbosity.


Example:

In your initialization code, you may set the log levels of debug messages:

   .
   .
   .

   reptilian.Goose  = goose.Alert(1) // Only error messages will be logged

   reptext.Goose    = goose.Alert(0) // No debug messages at all

   djparser.Goose   = goose.Alert(4) // Error messages and less verbose messages (levels 2~4) will be logged

   .
   .
   .


In the package to be debugged:

.
.
.

import (

   .
   .
   .

   "github.com/luisfurquim/goose"

   .
   .
   .

)

.
.
.

var Goose goose.Alert // Exported symbol needed only if you want to allow external control of the debug level

.
.
.


   // throughout your code you may log anything you want at varying levels of importance

   Goose.Logf(3, "Final Off=%d (%o)", d.Buf.Off, d.Buf.Off)

   .
   .
   .

   Goose.Logf(7, "Index: %#v", d.Index)

   // Logs will be actually printed only if the first parameter (the log level) is less equal than the log level indicated by the goose variable. Remember to never use the zero value, like Goose.Logf(0,...), as we want to make the log level 0 to print no debug messages at all.



