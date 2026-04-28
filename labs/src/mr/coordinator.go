package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
)

type Coordinator struct {
	mu         sync.Mutex
	files      []string
	nReduce    int   // fixed for this job: reduce task ids are 0 .. nReduce-1
	mapIdx     int   // next input file index to assign for map tasks
	reduceNext int   // next reduce task id to assign (starts at 0)
}

// Your code here -- RPC handlers for the worker to call.


//  Requesting a task from the RPC Handler, return a file with whether you are mapping or not
func (c *Coordinator) FindTask(args *TaskRequest, reply *TaskReply) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.mapIdx < len(c.files) {
		reply.Map = true
		reply.File = c.files[c.mapIdx]
		reply.MapTaskId = c.mapIdx
		reply.NReduce = c.nReduce
		c.mapIdx++
		return nil
	}
	if c.reduceNext < c.nReduce {
		reply.Map = false
		reply.ReduceId = int64(c.reduceNext)
		reply.NReduce = c.nReduce
		c.reduceNext++
		return nil
	}
	return nil
}

// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

// start a thread that listens for RPCs from worker.go
func (c *Coordinator) server(sockname string) {
	rpc.Register(c)
	rpc.HandleHTTP()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatalf("listen error %s: %v", sockname, e)
	}
	go http.Serve(l, nil)
}

// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.

	return ret
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(sockname string, files []string, nReduce int) *Coordinator {
	c := Coordinator{
		files:      files,
		nReduce:    nReduce,
		mapIdx:     0,
		reduceNext: 0,
	}

	c.server(sockname)
	return &c
}
