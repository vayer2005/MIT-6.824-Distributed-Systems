 
- Parallelism, fault tolerance, security kept in mind
- Partial failures are strange to deal with (network errors)

Implementation
- RPC (Remote Procedure Call)
- Threads/Concurrency Control

Performance
- Scalability - 2x resources -> 2x throughput
- Adding web servers that communicate with DB
	- DB bottlenecked

Fault Tolerance
- something breaks more often than for single computer systems
- Non volatile storage
- Replication -> issue is replicas may drift out of sync
	- Replicas should be as independent as possible

MapReduce (2004)
- Framework for distributing computation
- Chunked input files (web pages, files, etc.)
- Map function run on each input file (word to occurrence example)
	- all occurrences of a word aggregated over all maps per key call reduce
- Google File System
- 

