package main

import (
	"flag"
	"strings"

	"github.com/liamzebedee/pod-go/core"
)

func main() {
	replicasString := flag.String("replicas", "", "comma separated list of replicas")
	flag.Parse()

	replicaAddrs := []string{}
	if *replicasString != "" {
		replicaAddrs = strings.Split(*replicasString, ",")
	}

	replicaInfos := []core.ReplicaInfo{}
	for _, addr := range replicaAddrs {
		replicaInfos = append(replicaInfos, core.ReplicaInfo{DialAddress: addr})
	}

	client := core.NewClient()
	client.Start(replicaInfos)

	done := make(chan bool)
	<-done
}
