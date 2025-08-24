package consumer

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

type Message struct {
	ID int
}

func Handler(m *nats.Msg) {
	msg := &Message{}
	if err := json.Unmarshal(m.Data, msg); err != nil {
		log.Println(err)
	}
}
