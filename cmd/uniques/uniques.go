package main

import (
	"encoding/json"
	"log"
	"github.com/google/uuid"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// I'm not making two structs for this
type UniqueBody struct {
	maelstrom.MessageBody
	Id string `json:"id,omitempty"`
}


func main() {
	node := maelstrom.NewNode()
	node.Handle("generate", func(msg maelstrom.Message) error {
		var body UniqueBody;
		if err := json.Unmarshal(msg.Body,&body); err != nil {
			return err
		}
		body.Id = uuid.New().String()
		body.Type = "generate_ok";
		return node.Reply(msg,body)
	});
	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
