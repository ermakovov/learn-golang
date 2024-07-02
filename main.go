package main

import (
	"github.com/ermakovov/learn-golang/greeting"
	"github.com/ermakovov/learn-golang/webserver"
	"github.com/fatih/color"
)

func main() {
	color.Red(greeting.Hello())

	// webserver.StartCoursesWebserver()
	// webserver.StartMathWebserver()
	// webserver.StartCurrExchangeServer()
	// webserver.StartSocialNetworkServer()
	// webserver.StartArrayFinderServer()
	// webserver.StartSimpleStorageServer()
	webserver.StartURLExchangerServer()
}
