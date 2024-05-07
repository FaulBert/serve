package main

import (
	"flag"

	"github.com/nazhard/serve"
)

func main() {
	dir := flag.String("dir", ".", "set spesific directory to serve")
	port := flag.String("port", "9000", "set spesific port")
	log := flag.Bool("log", true, "doesn't print log if set to false")

	flag.Parse()

	server := &serve.New{
		Dir:  *dir,
		Port: *port,
		Log:  *log,
	}

	server.Start()
}
