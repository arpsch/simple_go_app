package main

import (
	"context"
	"flag"
	"log"
	"strings"

	"github.com/arpsch/ha/server"
)

var (
	port *string
	sid  *string
)

func init() {
	port = flag.String("port", ":8080", "a port for the server to listen, or defaults to :8080")
	sid = flag.String("sid", "",
		"sector id of this instance, or defaults to none i.e. the instance could used for multiple sectors")
}

func main() {
	err := doMain()
	if err != nil {
		log.Fatal(err)
	}
}

func doMain() error {

	ctx := context.Background()

	flag.Parse()
	// when ":" is not set, prefix it
	if idx := strings.Index(*port, ":"); idx == -1 {
		*port = ":" + *port
	}

	return server.InitAndRun(ctx, *port, *sid)
}
