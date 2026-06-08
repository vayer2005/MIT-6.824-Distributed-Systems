- Distributed transactions over multiple datacentres
- 2 phase commit over paxos participants
- Read only transactions are dominant
- Each shard has its own paxos protocol running. (replicated shards)
	- Minority replicas might be lagging (reading shards may be stale)
- Read Write Xaction
	- 1) Client picks unique txid to associate with each xaction
	- 2) Client sends reads to leaders of each paxos instance and each instance grabs lock and sends ack to client
	- 3) Client chooses one paxos group to be xaction cordinator. Client sends (Write request 
		- Xaction coordinator) to paxos leader
		- Paxos leader recieves write request and sends a prepare message to its followers (paxos log). Promises to carry out this transaction, majority response needed
		- Paxos leader sends yes to Xaction coordinator (another paxos leader)
	- 4) Xaction cooridinator requires YES from all shards needed for write
		- Xaction coordinator sends commit message and sends COMMIT REQUIRED to other leaders (once its own commit is saved in log)
	- 5) Every shard can execute write and release lock
	- 6) Why is Y the transaction leader?
- Read only Xaction
	- Serializable -> concurrency can be made linear
	- External consistency == Linearizability (must see most recent write) == Strong Consistency
	- What if we read the latest copy of local replica data
	- ![[Screenshot 2026-05-20 at 9.16.36 AM.png]]
	- Not serializable -> Reads split between transactions if you read latest value
	- Snapshot isolation (executes in timestamp order)
		- Note down each read timestamp -> x@10 = 348, x@20 = 443
		- R/w => commit time, R/0 => Start Time
		- Safe time
			- Delay until log record with nesessary timestamp has been recieved from leader
	- Issue -> Clock cannot be synchronized exactly
	- Commit wait rule
		- Waits until it knows its timestamp is in the past before committing
- Clock Synchronization
	- If clocks are not synced -> R/W transactions are fine because they use locks
		- R/O is too large -> waits for correct time
		- R/O is too small -> Miss recent committed writes (violates strong consistency)
	- True time scheme -> {Earliest, Latest}
	- Start Rule -> Xaction timestamp = latest of true time (guaranteed to have happened)
		- R/O -> Start, R/W Commit
	- 


 

Spanner Paper
- 
