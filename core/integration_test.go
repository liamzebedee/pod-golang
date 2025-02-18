package core

import (
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/liamzebedee/pod-go/core/pb"
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

func makeTx(data byte) *pb.Transaction {
	return &pb.Transaction{
		Ctx:   []byte{data},
		RMin:  0,
		RMax:  0,
		RConf: 0,
	}
}

func TestFlow(t *testing.T) {
	// Setup replicas config.
	replicas := make(map[PublicKey]*Replica)
	replicaConfigs := []ReplicaInfo{}
	N_REPLICAS := 5

	// Create replicas.
	for i := 0; i < N_REPLICAS; i++ {
		replica := NewReplica()
		addr := fmt.Sprintf("%s:%s", "127.0.0.1", getRandomPort())

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
	client.Write(makeTx(4))
	client.Write(makeTx(7))

	// Simulate for 5s then exit.
	ch := make(chan bool)
	go func() {
		time.Sleep(5 * time.Second)
		ch <- true
	}()
	<-ch
}
