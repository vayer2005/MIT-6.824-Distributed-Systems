
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
- 