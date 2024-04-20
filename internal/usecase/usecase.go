package usecase

import (
	"fmt"
	"log"
	"strings"

	"github.com/Nishad4140/SkillSync_ChatService/entities"
	"github.com/Nishad4140/SkillSync_ChatService/internal/adapter"
	"github.com/Nishad4140/SkillSync_ChatService/internal/usecase/chat"
)

type ChatUsecase struct {
	adapter adapter.ChatAdapterInterface
	Chat    *chat.ChatPool
}

func NewChatUsecase(adapter adapter.ChatAdapterInterface) ChatUsecase {
	return ChatUsecase{
		adapter: adapter,
		Chat:    chat.NewChatPool(),
	}
}

func (c *ChatUsecase) CreatePoolifnotalreadyExists(poolid string, insertChan chan<- entities.InsertIntoRoomMessage) (*chat.Pool, []entities.Message) {
	res, err := c.adapter.LoadMessages(poolid)
	ids := strings.Split(poolid, " ")
	if err != nil {
		log.Println("error while loading messages", err)
		res, err = c.adapter.LoadMessages(ids[1] + " " + ids[0])
		if err != nil {
			log.Println("error retrieving message ", err)
		}
	}
	if c.Chat.Pool[poolid] == nil {
		if c.Chat.Pool[ids[1]+" "+ids[0]] == nil {
			fmt.Println("no message found for id ", ids[1], ids[0])
			pool := chat.NewPool(ids[1] + " " + ids[0])
			go pool.Serve(insertChan)
			c.Chat.Pool[ids[1]+" "+ids[0]] = pool
			return pool, res
		}
		return c.Chat.Pool[ids[1]+" "+ids[0]], res

	}
	fmt.Println("pool id is ", poolid)
	return c.Chat.Pool[poolid], res
}
