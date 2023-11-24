package frameworks

import (
	"fmt"

	"github.com/google/uuid"
)

// This map can hold any type of data as the value
type Map map[string]interface{}

// Define VoteTypes
type VoteType int

// Type for scoring different votes
// High value for variables of this type expresses being in favour of vote.
type ScoreType float64

const (
	VoteToKickAgent VoteType = iota
	VoteToAcceptNewAgent
	VoteOnProposals
	VoteOnAllocation
)

type VoteInputs struct {
	DecisionType   VoteType    // Type of vote that needs to be made
	Candidates     []uuid.UUID // Map of choices [Dummy map for now]
	VoteParameters Map         // Parameters for the vote
}

type Vote struct {
	result map[uuid.UUID]interface{}
}

type VotingFramework struct {
	IDecisionFramework[VoteInputs, Vote]
}

func NewVotingFramework() *VotingFramework {
	return &VotingFramework{}
}

func (vf *VotingFramework) GetDecision(inputs VoteInputs) Vote {
	fmt.Println("VotingFramework: GetDecision called")
	fmt.Println("VotingFramework: Decision type: ", inputs.DecisionType)
	fmt.Println("VotingFramework: Choice map: ", inputs.Candidates)
	fmt.Println("VotingFramework: Vote parameters: ", inputs.VoteParameters)

	voteResult := vf.deliberateVote(inputs)

	return voteResult
}

func (vf *VotingFramework) deliberateVote(voteInputs VoteInputs) Vote {
	var vote Vote
	if voteInputs.DecisionType == VoteToKickAgent {
		// TODO: Deliberate on whether to kick an agent
		fmt.Println("Deliberating on whether to kick an agent")
		vote = VoteToKickWrapper(voteInputs)
	} else if voteInputs.DecisionType == VoteToAcceptNewAgent {
		// TODO: Deliberate on whether to accept a new agent
		fmt.Println("Deliberating on whether to accept a new agent")
	} else if voteInputs.DecisionType == VoteOnProposals {
		// TODO: Deliberate on how to vote on proposed directions
		fmt.Println("Deliberating on how to vote on proposals")
		//vote = Vote{result: Map{"decision": true}}
	} else if voteInputs.DecisionType == VoteOnAllocation {
		// TODO: Deliberate on how to vote on proposed directions
		fmt.Println("Deliberating on how to vote on proposals")
		vote = VoteOnAllocationWrapper(voteInputs)
	} else {
		// TODO: Deliberate on something else
		fmt.Println("Deliberating on something else")
		//vote = Vote{result: Map{"decision": true}}
	}
	return vote
}
