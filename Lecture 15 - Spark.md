- PageRank with Spark
	- First reads file (in reality will read per shard representation of file) -> collect()
	- What is a lineage graph? -> collect with execute across the lineage graph. Spark figures out where the data is and picks workers (compiles lineage graph)
	- map -> runs func on each element of input
		- spark spins off worker to each shard
	- join -> adds value to dict per key. dict.join(to_add)
	- Driver: runs the scala program
- HDFS split up all files
	- ![[Screenshot 2026-06-08 at 7.28.24 PM.png]]
	- Narrow stages then transform into shared wide stages
		- Bottlenecks over network (TB)?
			- Yes, this is heavyweight computation barrier
	- Failed worker in wide deps?
		- If worker fails and needs to restart, we need to restart entire computation because other workers have discarded their previous stages (needed for wide ops)
		- Spark allows checkpoints after specific transformations.

Spark paper
- RDDs -> resilient distributed datasets
	- Read only partitioned collection of records
	- Can run map, filter, or join transformations