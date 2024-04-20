package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Nishad4140/SkillSync_ChatService/entities"
	"github.com/Nishad4140/SkillSync_ChatService/internal/usecase"
	"github.com/Nishad4140/SkillSync_ChatService/internal/usecase/chat"
	"github.com/Nishad4140/SkillSync_ProtoFiles/pb"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

type ChatHandlers struct {
	InsertChannel chan<- entities.InsertIntoRoomMessage
	Usecase       usecase.ChatUsecaseInterface
	UserConn      pb.UserServiceClient
	Upgrader      websocket.Upgrader
}

func NewChatHandlers(insertChannel chan<- entities.InsertIntoRoomMessage, usecase usecase.ChatUsecaseInterface, userAddr string) *ChatHandlers {
	userRes, _ := grpc.Dial(userAddr, grpc.WithInsecure())
	return &ChatHandlers{
		InsertChannel: insertChannel,
		Usecase:       usecase,
		UserConn:      pb.NewUserServiceClient(userRes),
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (c *ChatHandlers) Handler(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("Sec-WebSocket-Extensions")
	userId := r.Header.Get("clientId")
	freelancerId := r.Header.Get("freelancerId")
	recieverId := r.Header.Get("recieverId")
	var poolId string
	if userId != "" && recieverId != "" {
		poolId = userId + " " + recieverId
	} else if freelancerId != "" && recieverId != "" {
		poolId = recieverId + " " + freelancerId
	} else {
		http.Error(w, "please provide valid headers", http.StatusBadRequest)
		return
	}
	clientData, err := c.UserConn.GetClientById(context.Background(), &pb.GetUserById{Id: userId})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name := clientData.Name
	clientId := userId
	if freelancerId != "" {
		freelancerData, err := c.UserConn.GetFreelancerById(context.Background(), &pb.GetUserById{
			Id: freelancerId,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		name = freelancerData.Name

	}
	conn, err := c.Upgrader.Upgrade(w, r, r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pool, msgs := c.Usecase.CreatePoolifnotalreadyExists(poolId, c.InsertChannel)
	client := chat.NewClient(conn, clientId, name, pool)
	client.Serve(msgs)
}

func (chat *ChatHandlers) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", chat.Handler)
	log.Println("listening on port 8000")
	if err := http.ListenAndServe(":8000", mux); err != nil {
		fmt.Println(err.Error())
	}
}
