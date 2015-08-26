package common

import (
	"time"
)


func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}


func Current() string {
	return time.Now().Format("2006-01-02 15:04:05.99999")
}

type RunError struct {
	error
}