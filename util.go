package main

import (
	"strconv"
	"time"
)

func genId() string {
	return strconv.Itoa(int(time.Now().UnixNano()))
}
