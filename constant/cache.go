package constant

import (
	"time"
)

type UserInfo struct {
	Username  string
	Ip        string
	Token     string
	Device    string
	LoginTime time.Time
}

const (
	LoggedInListKey string = "LoggedInList:%s:%s"
	TokenKey        string = "Token:%s"
	HeartbeatKey    string = "Heartbeat:%s"
)
