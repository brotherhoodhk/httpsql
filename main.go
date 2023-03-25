package main

import (
	"flag"
	"httpsql/controller"
)

var port = flag.Int("port", 8001, "port expose")
var extensionfunc = flag.Bool("extension", false, "is open extension function?")

func main() {
	flag.Parse()
	controller.ServerStart(*port)
}
