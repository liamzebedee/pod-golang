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
		Data: []byte{data},
	}
}

func TestFlow(t *testing.T) {
	fmt.Println("pod network test")
	fmt.Printf("system parameters: N=%d f=%d\n", pAlpha, pBeta)

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

	// Send transactions (stuttered).
	txs := []*pb.Transaction{
		makeTx(1),
		makeTx(2),
		makeTx(3),
		makeTx(4),
	}
	for _, tx := range txs {
		t.Logf("writing tx: %s", tx.ID())

		err := client.Write(tx)
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Simulate for a duration.
	SLEEP_TIME := 3 * time.Second
	ch := make(chan bool)
	go func() {
		time.Sleep(SLEEP_TIME)
		ch <- true
	}()

	<-ch

	// Read the pod and print the transaction timings.
	client.Read()
	t.Logf(client.ReadTx(txs[0].ID()).String())
}
