package adapter

import "github.com/Nishad4140/SkillSync_ChatService/entities"

type ChatAdapterInterface interface {
	InsertMessage(msg entities.InsertIntoRoomMessage) error
	LoadMessages(roomId string) ([]entities.Message, error)
}
