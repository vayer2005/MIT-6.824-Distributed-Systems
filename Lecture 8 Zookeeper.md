- If there exists a solution to the history that works then history is linearizable (graph of conds has no cycle)
- Each client sees the same history, only one sequence of ops
- No stale data; most recent completed write must be shown on reads
- Servers cache responses that the network may drop (up to design/implementation)
	- When client resends request sends old cached result.
- Zookeeper is a general purpose coordination service
- Can N x more servers lead to N x more results?
- Zoo runs on top of Zab (identical to raft for purpose of course)
	- More servers does not improve performance (leader is bottle-necked)
	- more servers slows performance as leader has to loop through more servers
- Client cannot send reads directly to clients as they may not be up to date with leader
- But zookeeper read performance goes up with more servers
	- non-linearizable reads. stale data allowed
	- Linearizable writes (raft?)

Zookeeper Paper Notes
- 