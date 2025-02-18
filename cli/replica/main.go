package main

import (
	"flag"
	"fmt"

	"github.com/liamzebedee/pod-go/core"
)

func main() {
	port := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	addr := fmt.Sprintf("127.0.0.1:%d", *port)

	replica := core.NewReplica()
	replica.ListenAndServe(addr)
}
