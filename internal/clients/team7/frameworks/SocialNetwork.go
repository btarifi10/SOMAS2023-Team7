package frameworks

import (
	"SOMAS2023/internal/common/utils"
	"fmt"

	"github.com/google/uuid"
)
const maxTrustIterations int = 6
type SocialConnection struct {
	connectionAge      int     // Number of rounds the agent has been known
	trustLevels  [3]float64 // Trust level of the agent, dummy float64 for now.
	isActiveConnection bool    // Boolean indicating if connection is on current bike
}

type SocialConnectionInput struct {
	agentDecisions map[uuid.UUID]utils.Forces
	agentDistribution ResourceAllocationParams
	agentLootBoxDecision utils.Colour
}

type ISocialNetwork[T any] interface {
	GetSocialNetwork() map[uuid.UUID]SocialConnection
	UpdateSocialNetwork(agentIds []uuid.UUID, inputs T)
	UpdateActiveConnections(agentIds []uuid.UUID)
	DeactivateConnections(agentIds []uuid.UUID)
	
}

type Itrust interface{
	GetAgentsOnCurrentBike() BaseBiker
	GetAgentByAgentId(agentId uuid.UUID) []float64
	GetVotemap() IdVoteMap 
}

type SocialNetwork struct {
	ISocialNetwork[SocialConnectionInput]
	Itrust
	socialNetwork *map[uuid.UUID]SocialConnection
}

func NewSocialNetwork() *SocialNetwork {
	return &SocialNetwork{
		socialNetwork: &map[uuid.UUID]SocialConnection{
			ConnectionAge: 0,
			trustLevels: [3]float64{0.5, 0.5, 0.5},
			isActiveConnection: false,
		},
	}
}

func (sn *SocialNetwork) GetSocialNetwork() map[uuid.UUID]SocialConnection {
	return *sn.socialNetwork
}

func (sn *SocialNetwork) updateTrustLevel(connection *SocialConnection, Input SocialConnectionInput) {
	
	// TODO: Update trust level based on forces
	agent:=env.GetAgentByAgentId
	DistancePenalty := sn.CalcDistributionPenalty(Input.agentDistribution)
	PedallingPenalty := sn.CalcPedallingPenalty(Input.forces.Pedal)
	OrientationPenalty := sn.CalcTurningPenalty(Input.forces.Turning)
	BrakingPenalty := sn.CalcBrakingPenalty(Input.forces.Brake)
	DifferentLootPenalty := sn.CalcDifferentLootBoxPenalty(Input.agentLootBoxDecision)

	W_dp := 1 
	W_op := 1
	W_bp := 1
	W_pp := 1
	W_dlp := 1
	
	// Calculate the new trust level
	trust := (W_dp*DistancePenalty) + (W_pp*PedallingPenalty) + (W_op*OrientationPenalty) + (W_bp*BrakingPenalty) + (W_dlp*DifferentLootPenalty)

    // Update the trust levels slice
    connection.TrustLevels = append(connection.TrustLevels[1:], trust)

    return trust
	}

func (sn *SocialNetwork) UpdateSocialNetwork(agentIds []uuid.UUID, inputs SocialConnectionInput) {
	for _, agentId := range agentIds {
		connection := (*sn.socialNetwork)[agentId]
		connection.connectionAge += 1
		sn.updateTrustLevel(&connection, inputs[agentId])
		(*sn.socialNetwork)[agentId] = connection
	}
}

func (sn *SocialNetwork) UpdateActiveConnections(agentIds []uuid.UUID) {
	for _, agentId := range agentIds {
		connection := (*sn.socialNetwork)[agentId]
		connection.isActiveConnection = true
		(*sn.socialNetwork)[agentId] = connection
	}
}

func (sn *SocialNetwork) DeactivateConnections(agentIds []uuid.UUID) {
	for _, agentId := range agentIds {
		connection := (*sn.socialNetwork)[agentId]
		connection.isActiveConnection = false
		(*sn.socialNetwork)[agentId] = connection
	}
}

// Retrieve agents on the current bike
func (sn *SocialNetwork) GetCurrentBikeNetwork() map[uuid.UUID]SocialConnection {
	activeConnections := map[uuid.UUID]SocialConnection{}
	for agentId, connection := range *sn.socialNetwork {
		if connection.isActiveConnection {
			activeConnections[agentId] = connection
		}
	}
	return activeConnections
}

// Implement individual calculation methods within the SocialNetwork

// Calc_Distribution_penalty calculates the penalty based on resources given
// and resources requested. This is a method of the SocialNetwork type.
func (sn *SocialNetwork) CalcDistributionPenalty(agentid uuid.UUID) float64 {

    return 1-GetVotemap()[agentid] // Return the calculated penalty
}



func (of *SocialNetwork) CalcPedallingPenalty(agentIds []uuid.UUID) float64 {
	return (1 - ((1 - agentIds.GetForces().Pedal) * agentIds.energyLevel)) + 0.1
}

func (of *SocialNetwork) CalcTurningPenalty(agentIds []uuid.UUID) float64 {
	if agentIds.GetForces().turning == of.Expected_turnangle() {
		return 1.1
	}
	return 0.7
}

func (of *SocialNetwork) CalcBrakingPenalty(forces utils.Forces) float64 {
	if sn.forces.brake > 0 {
		Points_brake = -1
		return Points_brake
	}
	else{
		Points_not_brake = 0
		return Points_not_brake	
	}
}

func (of *SocialNetwork) CalcDifferentLootBoxPenalty(agentIds []uuid.UUID) float64 {
	if of.GetColour() == agentIds.GetColour() {
		return 1.1
	}
	return 0.095
}

