package goose

import (
	"fmt"
	"log"
	"os"
)

type Alert uint8

func (d Alert) Logf(level int, format string, parms ...interface{}) {
	if uint8(d) >= uint8(level) {
		log.Printf(format, parms...)
	}
}

func (d Alert) Fatalf(level int, format string, parms ...interface{}) {
	if uint8(d) >= uint8(level) {
		log.Fatalf(format, parms...)
	}
	os.Exit(-1)
}

func (d Alert) Printf(level int, format string, parms ...interface{}) {
	if uint8(d) >= uint8(level) {
		fmt.Printf(format, parms...)
	}
}

func (d Alert) Sprintf(level int, format string, parms ...interface{}) string {
	if uint8(d) >= uint8(level) {
		return fmt.Sprintf(format, parms...)
	}
	return ""
}
