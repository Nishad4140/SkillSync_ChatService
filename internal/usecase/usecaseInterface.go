package usecase

import (
	"github.com/Nishad4140/SkillSync_ChatService/entities"
	"github.com/Nishad4140/SkillSync_ChatService/internal/usecase/chat"
)

type ChatUsecaseInterface interface {
	CreatePoolifnotalreadyExists(string, chan<- entities.InsertIntoRoomMessage) (*chat.Pool, []entities.Message)
}
