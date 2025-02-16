package core

type timestamp = float64

// Beta is the resilience threshold.
// Beta can be understood as the F in classical byzantine fault tolernace (BFT) systems.
var pBeta int = 1

// Alpha is the confirmation threshold.
// Pod-core security: alpha >= 4B+1
var pAlpha = 4*pBeta + 1

func GetParameters(numReplicas int, maxFailures int) (int, int) {
	minReplicas := 4*maxFailures + 1
	if numReplicas < minReplicas {
		panic("Number of replicas must be at least 4*maxFailures + 1")
	}

	alpha := minReplicas
	beta := maxFailures
	return alpha, beta
}
