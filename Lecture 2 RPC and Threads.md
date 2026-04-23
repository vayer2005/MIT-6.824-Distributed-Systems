
- https://go.dev/doc/effective_go
- Threads
	- Separate PC, Separate registers , separate stack.
	- I/O concurrency, Parallelism 
	- Convenience (polling goroutine to see if things have died)
- Event Driven 
	- Single threaded control
	- No CPU Paralellism
- Threads vs Processes
	- Inside a process you can have multiple threads
- Challenges with threading
	- Race conditions
- Coordination
	- Waitgroup, Cond variable

Web Crawler Multithreaded
- Cycles -> dont fetch pages twice
- Overlap network I/O
- run go program with -race flag to check if there are data races (not static analysis)

MapReduce
- abstraction that allows us to express the simple computations we were trying to perform but hides the messy details of parallelization, fault-tolerance, data distribution and load balancing in a library
- 



