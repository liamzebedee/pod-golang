package core

type Pod struct {
	// validators...
}

// – write(tx): It writes a transaction tx to the pod.
func (pod *Pod) Write(tx Transaction) {}

// – read() → D: It outputs a pod D = (T, rperf, Cpp).
func (pod *Pod) Read() {}
