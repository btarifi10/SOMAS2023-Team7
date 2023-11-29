package frameworks

import (
	"fmt"
)

// Ally needs the infrastructure to be updated for this to work
func VoteOnProposalsWrapper(voteInputs VoteInputs) Vote {
	var vote map[interface{}]interface{}
	switch voteInputs.VoteParameters {
	case Proportion:
		vote = Proportions(voteInputs)
	// TODO: Implement other voting possibilities
	//case YesNo:
	//vote = YesNos(voteInputs)
	default:
		fmt.Println("New decision type!")
		vote = Proportions(voteInputs)
	}

	return Vote{result: vote}
}

// TODO: Add functions for voting on which loot box to go to.

func Proportions(voteInputs VoteInputs) map[interface{}]interface{} {
	var candidates []interface{}
	var votes map[interface{}]interface{}
	candidates = voteInputs.Candidates
	totOptions := len(candidates)
	normalDist := 1.0 / float64(totOptions)
	//nearestLoot := GetNearestLootBox()
	for _, proposal_id := range candidates {
		/*
			// Basic implementation: only vote for nearest lootbox
			if proposal_id == GetNearestLootBox() {
				votes[proposal_id] = 1
			} else {
				votes[proposal_id] = 0
			}
		*/
		votes[proposal_id] = normalDist
	}
	return votes
}
