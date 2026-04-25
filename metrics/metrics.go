package metrics

import "sync/atomic"

var TotalRequests uint64
var FailedRequests uint64

func IncRequests(){
	atomic.AddUint64(&TotalRequests,1)
}

func IncFailures(){
	atomic.AddUint64(&FailedRequests,1)
}