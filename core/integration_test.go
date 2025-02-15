package core

import (
	"fmt"
	"net"
	"strconv"
	"testing"
)

func getRandomPort() string {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	portStr := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	fmt.Printf("got random port: %s\n", portStr)
	return portStr
}

func TestFlow(t *testing.T) {
	// Setup replicas config.
	replicas := make(map[PublicKey]*Replica)
	replicaConfigs := []ReplicaInfo{}
	N_REPLICAS := 5

	// Create replicas.
	for i := 0; i < N_REPLICAS; i++ {
		replica := NewReplica()
		addr := fmt.Sprintf("%s:%s", "localhost", getRandomPort())

		conf := ReplicaInfo{
			DialAddress: addr,
			PK:          replica.PublicKey(),
		}
		replicaConfigs = append(replicaConfigs, conf)
		replicas[replica.PublicKey()] = &replica

		go replica.ListenAndServe(conf.DialAddress)
	}

	// Create client.
	client := NewClient()

	// Start client with replica config.
	client.Start(replicaConfigs)

	// Send an example write.
	client.Write(makeTx(1))

	// Wait forever.
	ch := make(chan bool)
	<-ch
}
