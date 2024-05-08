package initializer

import (
	"github.com/Nishad4140/SkillSync_ChatService/internal/adapter"
	"github.com/Nishad4140/SkillSync_ChatService/internal/handler"
	"github.com/Nishad4140/SkillSync_ChatService/internal/usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

func Initializer(db *mongo.Database) *handler.ChatHandlers {
	adapter := adapter.NewChatAdapter(db)
	usecase := usecase.NewChatUsecase(adapter)
	insertRoom := usecase.InsertIntoDB()
	handler := handler.NewChatHandlers(insertRoom, &usecase, "ss-user-service:4001")
	return handler
}
