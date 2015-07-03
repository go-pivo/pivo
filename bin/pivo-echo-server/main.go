package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"gopkg.in/pivo.v2"
)

const MainBanner = `Pivo version %s (c) The Pivo Authors.`
const WelcomeBanner = `Welcome %s! Echo server is running Pivo %s`

var (
	lHttp = flag.String("http", ":8000", "listen for http on")
)

var server = NewServer()

func shutdown() {
	server.Stop()
	server.DisconnectAll()
	os.Exit(0)
}

func main() {
	fmt.Printf(MainBanner, pivo.Version)
	fmt.Println()
	flag.Parse()

	go server.Start()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; shutdown() }()

	// Listen for HTTP connections to upgrade
	if err := http.ListenAndServe(*lHttp, nil); err != nil {
		fmt.Println("websocket:", err)
		os.Exit(1)
	}
}
