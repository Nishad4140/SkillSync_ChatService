package entities

import "time"

type Message struct {
	Type    int
	Message string    `json:"Message" bson:"message"`
	Time    time.Time `json:"Time" bson:"time"`
	Name    string    `json:"Name" bson:"name"`
	Seen    bool      `json:"Seen" bson:"seen"`
}

type InsertIntoRoomMessage struct {
	RoomID   string  `bson:"room_id"`
	Messages Message `bson:"messages"`
}
