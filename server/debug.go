package main

import (
	"fmt"
)

// just a trick for now
func noDebug(s string, f ...interface{}) {}

func okDebug(s string, f ...interface{}) {
	fmt.Printf(s, f...)
}

var Debugf func(s string, f ...interface{}) = okDebug
