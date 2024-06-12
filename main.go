package main

import (
	"github.com/ermakovov/learn-golang/greeting"
	"github.com/fatih/color"
)

func main() {
	color.Red(greeting.Hello())
}
