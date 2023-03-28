package main

import (
	"flag"
	golangbackend "golang-backend"
)

func main() {
	addr := flag.String("addr", ":8080", "server addr")
	app := golangbackend.New(*addr)
	app.Start()
}
