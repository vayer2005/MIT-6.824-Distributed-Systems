package raft

// The file ../raftapi/raftapi.go defines the interface that raft must
// expose to servers (or the tester), but see comments below for each
// of these functions for more details.
//
// In addition,  Make() creates a new raft peer that implements the
// raft interface.

import (
	//	"bytes"
	"math/rand"
	"sync"
	"time"

	//	"6.5840/labgob"
	"6.5840/labrpc"
	"6.5840/raftapi"
	tester "6.5840/tester1"
)

// Ticker pacing: steady leader heartbeats (≤10/s) vs short follower/candidate polls.
const (
	heartbeatInterval    = 100 * time.Millisecond
	electionPollInterval = 10 * time.Millisecond
)

// A Go object implementing a single Raft peer.
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *tester.Persister   // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]

	// Your data here (3A, 3B, 3C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.

	currentTerm int
	votedFor    int
	log         []int
	commitIndex int
	lastApplied int

	serverState int // 0 follower, 1 leader, 2 candidate

	electionDeadline time.Time
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {

	var term int
	var isleader bool
	// Your code here (3A).
	rf.mu.Lock()
	defer rf.mu.Unlock()

	term = rf.currentTerm
	if rf.serverState == 1 {
		isleader = true
	} else {
		isleader = false
	}
	return term, isleader
}

// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
// before you've implemented snapshots, you should pass nil as the
// second argument to persister.Save().
// after you've implemented snapshots, pass the current snapshot
// (or nil if there's not yet a snapshot).
func (rf *Raft) persist() {
	// Your code here (3C).
	// Example:
	// w := new(bytes.Buffer)
	// e := labgob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// raftstate := w.Bytes()
	// rf.persister.Save(raftstate, nil)
}

// restore previously persisted state.
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	// Your code here (3C).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := labgob.NewDecoder(r)
	// var xxx
	// var yyy
	// if d.Decode(&xxx) != nil ||
	//    d.Decode(&yyy) != nil {
	//   error...
	// } else {
	//   rf.xxx = xxx
	//   rf.yyy = yyy
	// }
}

// how many bytes in Raft's persisted log?
func (rf *Raft) PersistBytes() int {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	return rf.persister.RaftStateSize()
}

// the service says it has created a snapshot that has
// all info up to and including index. this means the
// service no longer needs the log through (and including)
// that index. Raft should now trim its log as much as possible.
func (rf *Raft) Snapshot(index int, snapshot []byte) {
	// Your code here (3D).

}

// example RequestVote RPC arguments structure.
// field names must start with capital letters!
type RequestVoteArgs struct {
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

// example RequestVote RPC reply structure.
// field names must start with capital letters!
type RequestVoteReply struct {
	Term        int
	VoteGranted bool
}

// HEARTBEAT msg
// Sent from leader to follower, resets followers timeout. when recieved by follower
type ApppendEntriesArgs struct {
}

type ApppendEntriesReply struct {
}

func randomElectionTimeout() int64 {
	return 250 + (rand.Int63() % 300)
}

func (rf *Raft) AppendEntries(args *ApppendEntriesArgs, reply *ApplyErrReply) {
	//TODO: reset timeout var.
	rf.mu.Lock()
	defer rf.mu.Unlock()
	rf.electionDeadline = time.Now().Add(time.Duration(randomElectionTimeout()))

}

// example RequestVote RPC handler.
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	reply.Term = rf.currentTerm
	reply.VoteGranted = false

	// Stale term — reject (Figure 2).
	if args.Term < rf.currentTerm {
		return
	}

	// At least as new as us then adopt term and become follower if we were candidate/leader.
	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		rf.votedFor = -1
		rf.serverState = 0
	}

	lastIdx := len(rf.log) - 1
	lastTerm := 0
	if lastIdx >= 0 {
		lastTerm = rf.log[lastIdx]
	}
	upToDate := args.LastLogTerm > lastTerm ||
		(args.LastLogTerm == lastTerm && args.LastLogIndex >= lastIdx)

	if (rf.votedFor == -1 || rf.votedFor == args.CandidateId) && upToDate {
		rf.votedFor = args.CandidateId
		reply.VoteGranted = true
		rf.electionDeadline = time.Now().Add(time.Duration(randomElectionTimeout()) * time.Millisecond)
	}

	reply.Term = rf.currentTerm
}

// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true

	// Your code here (3B).

	return index, term, isLeader
}

func (rf *Raft) startElection() {
	//TODO; already holding lock

	rf.currentTerm += 1
	lastIdx := len(rf.log) - 1
	lastTerm := 0
	if lastIdx >= 0 {
		lastTerm = rf.log[lastIdx]
	}
	args := RequestVoteArgs{
		Term:         rf.currentTerm,
		CandidateId:  rf.me,
		LastLogIndex: lastIdx,
		LastLogTerm:  lastTerm,
	}

	for i := range rf.peers {
		if i != rf.me {
			reply := RequestVoteReply{}
			_ = rf.peers[i].Call("Raft.RequestVote", &args, &reply)
		}
	}


}

func (rf *Raft) checkElect() {
	//TODO: periodically check for re-election

	for {
		rf.mu.Lock()
		if time.Now().After(rf.electionDeadline) {
			//Kick off re-election
			rf.startElection()
		}
		rf.mu.Unlock()
		time.Sleep(electionPollInterval)
		
	}

}



func (rf *Raft) tickFollower() {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	// Bootstrap deadline once so After(deadline) is meaningful before AppendEntries resets it.
	if rf.electionDeadline.IsZero() {
		rf.electionDeadline = time.Now().Add(time.Duration(randomElectionTimeout()) * time.Millisecond)
		return
	}

	// Follower (0) or candidate (2): election watchdog only — do not bump deadline here.
	if (rf.serverState == 0 || rf.serverState == 2) && time.Now().After(rf.electionDeadline) {
		rf.startElection()
	}
}

func (rf *Raft) ticker() {
	for {
		rf.mu.Lock()
		st := rf.serverState
		rf.mu.Unlock()

		if st == 1 {
			for i := range rf.peers {
				if i != rf.me {
					args := &ApppendEntriesArgs{}
					reply := &ApppendEntriesReply{}
					_ = rf.peers[i].Call("Raft.AppendEntries", args, reply)
				}
			}
			time.Sleep(heartbeatInterval)
		} else {
			rf.tickFollower()
			time.Sleep(electionPollInterval)
		}
	}
}

// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
func Make(peers []*labrpc.ClientEnd, me int,
	persister *tester.Persister, applyCh chan raftapi.ApplyMsg) raftapi.Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me

	// Your initialization code here (3A, 3B, 3C).
	rf.commitIndex = 0
	rf.lastApplied = 0

	rf.serverState = 0 // start as follower
	rf.votedFor = -1
	rf.log = []int{0} // dummy entry at index 0 (term 0); avoids empty-log edge cases

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	// start ticker goroutine to start elections
	go rf.ticker()

	// election checker process

	go rf.checkElect()

	return rf
}
