package mr

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/rpc"
	"os"
	"path/filepath"
	"sort"
)


// Map functions return a slice of KeyValue.
type KeyValue struct {
	Key   string
	Value string
}

// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

var coordSockName string // socket for coordinator

// encodeToFile appends one JSON-encoded KeyValue to mr-<mapTask>-<reduceId>.
func encodeToFile(reduceId int, mapTask int, kv KeyValue) error {
	name := fmt.Sprintf("mr-%d-%d", mapTask, reduceId)
	f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(&kv)
}


// decodeFromFilesAndReduce reads all mr-*-<reduceId> shards (same layout as mrsequential after merge).
func decodeFromFilesAndReduce(reduceId int, reducef func(string, []string) string) error {
	pattern := fmt.Sprintf("mr-*-%d", reduceId)
	filenames, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	var kva []KeyValue
	for _, name := range filenames {
		f, err := os.Open(name)
		if err != nil {
			return err
		}
		dec := json.NewDecoder(f)
		for {
			var kv KeyValue
			if err := dec.Decode(&kv); err != nil {
				if err == io.EOF {
					break
				}
				f.Close()
				return err
			}
			kva = append(kva, kv)
		}
		f.Close()
	}

	sort.Slice(kva, func(i, j int) bool {
		return kva[i].Key < kva[j].Key
	})

	oname := fmt.Sprintf("mr-out-%d", reduceId)
	ofile, err := os.Create(oname)
	if err != nil {
		return err
	}
	defer ofile.Close()

	i := 0
	for i < len(kva) {
		j := i + 1
		for j < len(kva) && kva[j].Key == kva[i].Key {
			j++
		}
		values := []string{}
		for k := i; k < j; k++ {
			values = append(values, kva[k].Value)
		}
		output := reducef(kva[i].Key, values)
		fmt.Fprintf(ofile, "%v %v\n", kva[i].Key, output)
		i = j
	}

	return nil
}

// main/mrworker.go calls this function.
func Worker(sockname string, mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	coordSockName = sockname

	//Call cordinator to get a task
	args := TaskRequest{}
	reply := TaskReply{}

	ok := call("Coordinator.FindTask", &args, &reply)
	
	if ok && reply.Map {
		data, err := os.ReadFile(reply.File)
		if err != nil {
			fmt.Print(err)
			return
		}
		contents := string(data)
		kva := mapf(reply.File, contents)
		mapTask := reply.MapTaskId
		nReduce := reply.NReduce

		for i := range kva {
			kv := kva[i]
			reduceId := ihash(kv.Key) % nReduce

			if err := encodeToFile(reduceId, mapTask, kv); err != nil {
				log.Printf("encodeToFile %v", err)
				return
			}
		}
		
	} else if ok {
		if err := decodeFromFilesAndReduce(int(reply.ReduceId), reducef); err != nil {
			log.Printf("decodeFromFilesAndReduce: %v", err)
			return
		}

	} else {
		fmt.Printf("call failed!\n")
	}
	

	// uncomment to send the Example RPC to the coordinator.
	// CallExample()

}

// example function to show how to make an RPC call to the coordinator.
//
// the RPC argument and reply types are defined in rpc.go.
func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply.
	// the "Coordinator.Example" tells the
	// receiving server that we'd like to call
	// the Example() method of struct Coordinator.
	ok := call("Coordinator.Example", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("reply.Y %v\n", reply.Y)
	} else {
		fmt.Printf("call failed!\n")
	}
}

// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	c, err := rpc.DialHTTP("unix", coordSockName)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	if err := c.Call(rpcname, args, reply); err == nil {
		return true
	}
	log.Printf("%d: call failed err %v", os.Getpid(), err)
	return false
}
