package goose
import (
    "log/syslog"
    "testing"
)


func TestSyslog(t *testing.T) {
	err := UseSyslogNet("tcp", "127.0.0.1:514", syslog.LOG_EMERG)
	if err != nil {
		t.Errorf("Failed UseSyslogNet: %s", err)
	}
}

