package frameworks

import (
	"SOMAS2023/internal/common/utils"
	"math"

	"github.com/google/uuid"
)

const maxTrustIterations int = 6

type SocialConnection struct {
	connectionAge      int                         // Number of rounds the agent has been known
	trustLevels        [maxTrustIterations]float64 // Trust level of the agent, dummy float64 for now.
	isActiveConnection bool                        // Boolean indicating if connection is on current bike
}

type SocialConnectionInput struct {
	agentDecisions       map[uuid.UUID]utils.Forces
	agentDistribution    ResourceAllocationParams
	agentLootBoxDecision utils.Colour
}

type ISocialNetwork[T any] interface {
	GetSocialNetwork() map[uuid.UUID]SocialConnection
	UpdateSocialNetwork(agentIds []uuid.UUID, inputs T)
	UpdateActiveConnections(agentIds []uuid.UUID)
	DeactivateConnections(agentIds []uuid.UUID)
}

type Itrust interface {
	GetAgentsOnCurrentBike() BaseBiker
	GetAgentByAgentId(agentId uuid.UUID) IBaseBiker
	GetVotemap(agentId uuid.UUID) IdVoteMap //get votes of each individual agent
	GetcallingAgentId() uuid.UUID
	GetleaderID() uuid.UUID
	GetProposedTurnangle() utils.TurningDecision
	GetAgentscolorByagentID(agentId uuid.UUID) utils.Colour
}

type SocialNetwork struct {
	ISocialNetwork[SocialConnectionInput]
	Itrust
	socialNetwork *map[uuid.UUID]SocialConnection
	personality   *Personality
}

func NewSocialNetwork(p *Personality) *SocialNetwork {
	return &SocialNetwork{
		socialNetwork: &map[uuid.UUID]SocialConnection{},
		personality:   p,
	}
}

func (sn *SocialNetwork) GetSocialNetwork() map[uuid.UUID]SocialConnection {
	return *sn.socialNetwork
}

func (sn *SocialNetwork) UpdateTrustLevel(agentId uuid.UUID, connection *SocialConnection, Input SocialConnectionInput, p *Personality) {

	// TODO: Update trust level based on forces
	agent := sn.GetAgentByAgentId(agentId)
	DistancePenalty := sn.CalcDistributionPenalty(agentId, sn.Itrust.GetVotemap(agentId), sn.Itrust.GetAgentsOnCurrentBike(), p)
	PedallingPenalty := sn.CalcPedallingPenalty(agentId, Input[agentId].Pedal, p, sn.Itrust.GetAgentsOnCurrentBike())
	OrientationPenalty := sn.CalcTurningPenalty(agentId)
	BrakingPenalty := sn.CalcBrakingPenalty(Input[agentId].Brake)
	DifferentLootPenalty := sn.CalcDifferentLootBoxPenalty(Input.agentLootBoxDecision, p)

	W_dp := 1
	W_op := 1
	W_bp := 1
	W_pp := 1
	W_dlp := 1

	// Calculate the new trust level
	trust := (W_dp * DistancePenalty) + (W_pp * PedallingPenalty) + (W_op * OrientationPenalty) + (W_bp * BrakingPenalty) + (W_dlp * DifferentLootPenalty)

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
func (sn *SocialNetwork) CalcDistributionPenalty(agentId uuid.UUID, resourcedistribution IdVoteMap, bikers []BaseBiker, p *Personality) float64 {
	var penalty float64
	switch {
	case p.Egalitarian:
		penalty := 0.0 // Assuming penalty is a float64
		bikerCount := len(bikers)
		for _, agent := range bikers {

			expectedValue := 1.0 / float64(bikerCount)
			vote := resourcedistribution[agent.GetID()]
			penalty += math.Abs(expectedValue - vote)
		}
	case p.Selfish:
		for _, agent := range bikers {
			if agent.GetID() == sn.GetcallingAgentId() {
				penalty = 1 - resourcedistribution[agent.GetID()]
			}
		}
	case p.Judgemental:
		for _, agent := range bikers {
			if agent.GetID() == agentId {
				penalty = resourcedistribution[agent.GetID()]
			}
		}

	case p.Utilitarian:
		for _, agent := range bikers {

			vote := resourcedistribution[agent.GetID()]
			penalty += math.Abs(vote - (1 - agent.GetEnergyLevel()))
		}
	}

	return penalty // Return the calculated penalty
}

//Find shift to account for forgiveness

func (sn *SocialNetwork) CalcPedallingPenalty(agentId uuid.UUID, pedalling_force float64, p *Personality, bikers []BaseBiker) float64 {
	var penalty float64
	switch {
	case p.Egalitarian:
		Total_pedalling := 0.0 // Assuming penalty is a float64
		bikerCount := len(bikers)
		for _, agent := range bikers {
			Total_pedalling += agent.GetForces().Pedal
		}
		expectedValue := Total_pedalling / float64(bikerCount)

		penalty = expectedValue - pedalling_force
		return penalty

	case p.Utilitarian:
		for _, agent := range bikers {

			penalty = ((1 - pedalling_force) * sn.GetAgentByAgentId(agentId).GetEnergyLevel()) - (math.Pow(pedalling_force, 1.0/pedalling_force))*0.3
			return penalty
		}
	}

	return 1 - pedalling_force // Return the calculated penalty
}

func (sn *SocialNetwork) CalcTurningPenalty(agentId uuid.UUID) float64 {
	if agentId == sn.GetleaderID() {
		if sn.GetAgentByAgentId(agentId).GetForces.Turning == sn.GetProposedTurnangle() {
			return -0.2
		}
		return 0.3
	}
	if sn.GetAgentByAgentId(agentId).GetForces.Turning == 0 {
		return -0.2
	}
	return 0.5
}

func (sn *SocialNetwork) CalcBrakingPenalty(braking float64) float64 {
	return 0.8
}

func (sn *SocialNetwork) CalcDifferentLootBoxPenalty(agentIds []uuid.UUID) float64 {
	if sn.GetAgentscolorByagentID(sn.GetcallingAgentId()) == sn.GetAgentscolorByagentID(sn.GetcallingAgentId()) {
		return -0.2
	}
	return 0.095
}
