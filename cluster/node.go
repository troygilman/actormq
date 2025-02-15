package cluster

import (
	"errors"
	"log/slog"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/troygilman/actormq/cluster/timer"
)

type (
	heartbeatTimeout struct{}
	electionTimeout  struct{}
)

type NodeConfig struct {
	Topic               string
	DiscoveryPID        *actor.PID
	StateMachine        *actor.PID
	Logger              *slog.Logger
	ElectionMinServers  uint64
	ElectionMinInterval time.Duration
	ElectionMaxInterval time.Duration
	HeartbeatInterval   time.Duration
}

func NewNodeConfig() NodeConfig {
	return NodeConfig{
		ElectionMinServers:  3,
		ElectionMinInterval: 150 * time.Millisecond,
		ElectionMaxInterval: 300 * time.Millisecond,
		HeartbeatInterval:   50 * time.Millisecond,
	}
}

func (config NodeConfig) WithDiscoveryPID(discoveryPID *actor.PID) NodeConfig {
	config.DiscoveryPID = discoveryPID
	return config
}

func (config NodeConfig) WithLogger(logger *slog.Logger) NodeConfig {
	config.Logger = logger
	return config
}

type nodeActor struct {
	config            NodeConfig
	leader            *actor.PID
	currentTerm       uint64
	votedFor          *actor.PID
	log               []*LogEntry
	commitIndex       uint64
	lastApplied       uint64
	votes             uint64
	nodes             map[uint64]*nodeMetadata
	pendingCommands   map[uint64]*commandMetadata
	heartbeatRepeater actor.SendRepeater
	electionTimer     *timer.SendTimer
}

func NewNode(config NodeConfig) actor.Producer {
	return func() actor.Receiver {
		return &nodeActor{
			config: config,
		}
	}
}

func (node *nodeActor) Receive(act *actor.Context) {
	switch msg := act.Message().(type) {
	case actor.Initialized:
		node.nodes = make(map[uint64]*nodeMetadata)
		node.pendingCommands = make(map[uint64]*commandMetadata)

	case actor.Started:
		node.electionTimer = timer.NewSendTimer(act.Engine(), act.PID(), electionTimeout{}, newElectionTimoutDuration(node.config))
		node.heartbeatRepeater = act.SendRepeat(act.PID(), heartbeatTimeout{}, node.config.HeartbeatInterval)
		act.Send(node.config.DiscoveryPID, &RegisterNode{Topic: node.config.Topic})

	case *ActiveNodes:
		node.handleActiveNodes(act, msg)

	case *actor.Ping:
		act.Send(act.Sender(), &actor.Pong{})

	case *Envelope:
		node.handleEnvelope(act, msg)

	case *AppendEntries:
		node.handleExternalTerm(msg.Term)
		node.handleAppendEntries(act, msg)

	case *AppendEntriesResult:
		node.handleExternalTerm(msg.Term)
		node.handleAppendEntriesResult(act, msg)

	case *RequestVote:
		node.handleExternalTerm(msg.Term)
		node.handleRequestVote(act, msg)

	case *RequestVoteResult:
		node.handleExternalTerm(msg.Term)
		node.handleRequestVoteResult(act, msg)

	case electionTimeout:
		node.electionTimer.Reset(newElectionTimoutDuration(node.config))
		if !pidEquals(act.PID(), node.leader) {
			node.startElection(act)
		}

	case heartbeatTimeout:
		if pidEquals(node.leader, act.PID()) {
			node.sendAppendEntriesAll(act)
		}
	}

	node.updateStateMachine(act)
}

func (node *nodeActor) handleActiveNodes(act *actor.Context, msg *ActiveNodes) {
	node.nodes = make(map[uint64]*nodeMetadata)
	lastLogIndex, _ := node.lastLogIndexAndTerm()
	for _, pid := range msg.Nodes {
		pid := PIDToActorPID(pid)
		if !pidEquals(pid, act.PID()) {
			key := pid.LookupKey()
			if _, ok := node.nodes[key]; !ok {
				node.nodes[key] = &nodeMetadata{
					pid:        pid,
					nextIndex:  lastLogIndex + 1,
					matchIndex: 0,
				}
			}
		}
	}
	node.config.Logger.Info("handleActiveNodes", "msg", msg, "nodes", node.nodes)
}

func (node *nodeActor) handleEnvelope(act *actor.Context, msg *Envelope) {
	node.config.Logger.Info("handleMessage", "pid", act.PID(), "sender", act.Sender(), "msg", msg)
	if pidEquals(node.leader, act.PID()) {
		node.log = append(node.log, &LogEntry{
			Message: msg.Message,
			Term:    node.currentTerm,
		})
		newLogIndex := uint64(len(node.log))
		node.pendingCommands[newLogIndex] = &commandMetadata{
			sender: act.Sender(),
		}
		node.sendAppendEntriesAll(act)
	} else {
		var redirectPID *PID
		if node.leader != nil {
			redirectPID = ActorPIDToPID(node.leader)
		}
		act.Send(act.Sender(), &EnvelopeResult{
			Success:     false,
			RedirectPID: redirectPID,
		})
	}
}

func (node *nodeActor) handleAppendEntries(act *actor.Context, msg *AppendEntries) {
	result := &AppendEntriesResult{}
	defer func() {
		result.Term = node.currentTerm
		act.Send(act.Sender(), result)
		node.config.Logger.Debug("handleAppendEntries", "pid", act.PID(), "sender", act.Sender(), "msg", msg, "result", result)
	}()

	if msg.Term == node.currentTerm && pidEquals(node.leader, act.PID()) {
		node.config.Logger.Warn("Leader collision!")
	}

	// Condition #1
	// Reply false if term < currentTerm
	if msg.Term < node.currentTerm {
		result.Success = false
		return
	}

	node.leader = act.Sender()

	// Condition #2
	// Reply false if log doesn't contain an entry at prevLogIndex whose term matches prevLogTerm
	if msg.PrevLogIndex > 0 {
		if len(node.log) < int(msg.PrevLogIndex) || (len(node.log) > 0 && node.log[msg.PrevLogIndex-1].Term != msg.PrevLogTerm) {
			result.Success = false
			return
		}
	}

	newEntryIndex := msg.PrevLogIndex
	for _, entry := range msg.Entries {
		newEntryIndex++

		// Condition #3
		// If an existing entry conflicts with a new one (same index but different terms),
		// delete the existing entry and all that follow it
		if len(node.log) >= int(newEntryIndex) && node.log[newEntryIndex-1].Term != entry.Term {
			node.log = node.log[:newEntryIndex-1]
		}

		// Condition #4
		// Append any new entries not already in the log
		if len(node.log) < int(newEntryIndex) {
			node.log = append(node.log, entry)
		}
	}

	// Condition #5
	// If leaderCommit > commitIndex,
	// set commitIndex = min(leaderCommit, index of last new entry)
	if msg.LeaderCommit > node.commitIndex {
		node.commitIndex = min(msg.LeaderCommit, newEntryIndex)
	}

	result.Success = true

	node.electionTimer.Reset(newElectionTimoutDuration(node.config))
}

func (node *nodeActor) handleAppendEntriesResult(act *actor.Context, msg *AppendEntriesResult) {
	node.config.Logger.Debug("handleAppendEntriesResult", "pid", act.PID(), "sender", act.Sender(), "msg", msg)
	metadata, ok := node.nodes[act.Sender().LookupKey()]
	if !ok {
		node.config.Logger.Error("handleAppendEntriesResult", "pid", act.PID(), "sender", act.Sender(), "msg", msg, "error", errors.New("could not find PID"))
		return
	}
	if msg.Success {
		lastLogIndex, _ := node.lastLogIndexAndTerm()
		metadata.matchIndex = metadata.nextIndex
		metadata.nextIndex = lastLogIndex + 1
	} else {
		if metadata.nextIndex > 1 {
			metadata.nextIndex--
		}
		if err := node.sendAppendEntries(act, metadata.pid); err != nil {
			node.config.Logger.Error("handleAppendEntriesResult", "pid", act.PID(), "sender", act.Sender(), "msg", msg, "error", err)
		}
	}
}

func (node *nodeActor) handleRequestVote(act *actor.Context, msg *RequestVote) {
	result := &RequestVoteResult{}
	defer func() {
		result.Term = node.currentTerm
		act.Send(act.Sender(), result)
		node.config.Logger.Info("handleRequestVote", "pid", act.PID(), "sender", act.Sender(), "msg", msg, "result", result)
	}()

	// Condition #1
	// Reply false if term < currentTerm
	if msg.Term < node.currentTerm {
		result.VoteGranted = false
		return
	}

	// Condition #2
	// If votedFor is null or candidateId,
	// and candidate's log is at least as up-to-date as receiver's log, grant vote
	candidatePID := act.Sender()
	if node.votedFor == nil || node.votedFor.String() == candidatePID.String() {
		if msg.LastLogIndex >= node.lastApplied && (node.lastApplied == 0 || msg.LastLogTerm >= node.log[node.lastApplied-1].Term) {
			node.votedFor = candidatePID
			result.VoteGranted = true
		}
	}
}

func (node *nodeActor) handleRequestVoteResult(act *actor.Context, msg *RequestVoteResult) {
	node.config.Logger.Info("handleRequestVoteResult", "pid", act.PID(), "sender", act.Sender(), "msg", msg)
	if msg.VoteGranted && msg.Term == node.currentTerm && !pidEquals(node.leader, act.PID()) {
		node.votes++
		if float32(node.votes)/float32(len(node.nodes)) > 0.5 {
			node.config.Logger.Info("Promoted to leader", "pid", act.PID(), "sender", act.Sender())
			node.leader = act.PID()
			lastLogIndex, _ := node.lastLogIndexAndTerm()
			for _, metadata := range node.nodes {
				metadata.nextIndex = lastLogIndex + 1
				metadata.matchIndex = 0
			}
			node.sendAppendEntriesAll(act)
		}
	}
}

func (node *nodeActor) sendAppendEntriesAll(act *actor.Context) {
	for _, metadata := range node.nodes {
		if err := node.sendAppendEntries(act, metadata.pid); err != nil {
			node.config.Logger.Error("Sending AppendEntries for "+metadata.pid.String(), "pid", act.PID(), "error", err.Error())
		}
	}
}

func (node *nodeActor) sendAppendEntries(act *actor.Context, pid *actor.PID) error {
	metadata, ok := node.nodes[pid.LookupKey()]
	if !ok {
		return errors.New("server does not exist")
	}

	if metadata.nextIndex == 0 {
		return errors.New("nextIndex is 0 for " + pid.String())
	}

	entries := []*LogEntry{}
	lastLogIndex, _ := node.lastLogIndexAndTerm()
	if lastLogIndex >= metadata.nextIndex {
		entries = node.log[metadata.nextIndex-1:]
	}

	var prevLogIndex uint64 = metadata.nextIndex - 1
	var prevLogTerm uint64 = 0
	if prevLogIndex > 0 {
		prevLogTerm = node.log[prevLogIndex-1].Term
	}

	act.Send(metadata.pid, &AppendEntries{
		Term:         node.currentTerm,
		PrevLogTerm:  prevLogTerm,
		PrevLogIndex: prevLogIndex,
		Entries:      entries,
		LeaderCommit: node.commitIndex,
	})
	return nil
}

func (node *nodeActor) startElection(act *actor.Context) {
	defer func() {
		node.config.Logger.Info("Starting election", "pid", act.PID(), "term", node.currentTerm)
	}()
	node.currentTerm++
	node.votes = 1
	node.votedFor = act.PID()

	if len(node.nodes)+1 < int(node.config.ElectionMinServers) {
		node.config.Logger.Warn("Not enough servers for election", "pid", act.PID())
		return
	}

	lastLogIndex, lastLogTerm := node.lastLogIndexAndTerm()
	for _, metadata := range node.nodes {
		act.Send(metadata.pid, &RequestVote{
			Term:         node.currentTerm,
			LastLogIndex: lastLogIndex,
			LastLogTerm:  lastLogTerm,
		})
	}
}

func (node *nodeActor) lastLogIndexAndTerm() (uint64, uint64) {
	var lastLogIndex uint64 = uint64(len(node.log))
	var lastLogTerm uint64 = 0
	if lastLogIndex > 0 {
		lastLogTerm = node.log[lastLogIndex-1].Term
	}
	return lastLogIndex, lastLogTerm
}

func (node *nodeActor) handleExternalTerm(term uint64) {
	if term > node.currentTerm {
		node.currentTerm = term
		node.leader = nil
		node.votedFor = nil
	}
}

func (node *nodeActor) updateStateMachine(act *actor.Context) {
	if pidEquals(node.leader, act.PID()) {
		for i := uint64(len(node.log)); i >= node.commitIndex+1; i-- {
			if node.log[i-1].Term == node.currentTerm {
				matched := 0
				for _, metadata := range node.nodes {
					if metadata.matchIndex >= i {
						matched++
					}
				}
				if float32(matched) > float32(len(node.nodes))/2 {
					node.commitIndex = i
					break
				}
			}
		}
	}
	for node.commitIndex > node.lastApplied {
		node.lastApplied++
		entry := node.log[node.lastApplied-1]
		node.applyMessage(act, entry.GetMessage())
		command, ok := node.pendingCommands[node.lastApplied]
		if ok {
			act.Send(command.sender, &EnvelopeResult{
				Success: true,
			})
			delete(node.pendingCommands, node.lastApplied)
		}
		node.config.Logger.Info("Applied message", "pid", act.PID(), "index", node.lastApplied, "msg", entry.Message)
	}
}

func (node *nodeActor) applyMessage(act *actor.Context, msg *Message) {
	_, err := act.Request(node.config.StateMachine, &ConsumerEnvelope{
		Message:  msg,
		IsLeader: pidEquals(node.leader, act.PID()),
	}, time.Second).Result()
	if err != nil {
		panic(err)
	}
}
