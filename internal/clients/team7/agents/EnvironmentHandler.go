package agents

import (
	objects "SOMAS2023/internal/common/objects"
	utils "SOMAS2023/internal/common/utils"

	"github.com/google/uuid"
)

type EnvironmentHandler struct {
	GameState     objects.IGameState // Game state to be updated in each round
	CurrentBikeId uuid.UUID          // unique ID of current bike
}

func (env *EnvironmentHandler) GetLootBoxesByColour(colour utils.Colour) []objects.ILootBox {
	lootBoxes := env.GameState.GetLootBoxes()
	var matchingLootBoxes []objects.ILootBox
	for _, lootBox := range lootBoxes {
		if lootBox.GetColour() == colour {
			matchingLootBoxes = append(matchingLootBoxes, lootBox)
		}
	}
	return matchingLootBoxes
}

func (env *EnvironmentHandler) GetAgentsOnCurrentBike() []objects.IBaseBiker {
	return env.GetBikeAgentsByBikeId(env.CurrentBikeId)
}

func (env *EnvironmentHandler) GetBikeAgentsByBikeId(bikeId uuid.UUID) []objects.IBaseBiker {
	megaBikes := env.GameState.GetMegaBikes()
	bike := megaBikes[bikeId]
	return bike.GetAgents()
}

func (env *EnvironmentHandler) GetBikerListByAgentIds(agentIds []uuid.UUID) []objects.IBaseBiker {
	var bikers []objects.IBaseBiker
	for _, bike := range env.GameState.GetMegaBikes() {
		for _, agent := range bike.GetAgents() {
			for _, agentId := range agentIds {
				if agentId == agent.GetID() {
					bikers = append(bikers, agent)
				}
			}
		}
	}
	return bikers
}

func NewEnvironmentHandler(gameState objects.IGameState, bikeId uuid.UUID) *EnvironmentHandler {
	return &EnvironmentHandler{
		GameState:     gameState,
		CurrentBikeId: bikeId,
	}
}
