package example

import "time"

const (
	Addr = "127.0.0.1:8080"
	Path = "/state"
)

type ReportREQ struct {
	State string
}

type StateACK struct {
	State string
	Time  time.Time
}
