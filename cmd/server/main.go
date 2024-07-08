package main

import (
	"flag"
	golangbackend "golang-backend"
)

func main() {
	addr := flag.String("addr", ":8080", "server addr")
	rollback := flag.Bool("rollback", false, "rollback migration")
	app := golangbackend.New(*addr, true, true, *rollback)
	app.Start()
}
