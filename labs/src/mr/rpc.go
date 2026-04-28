package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

// Add your RPC definitions here.

type TaskRequest struct{}

type TaskReply struct {
	Map       bool // True if current task is a Map
	File      string
	ReduceId  int64
	NReduce   int // exported for RPC (worker uses ihash % NReduce)
	MapTaskId int
}

