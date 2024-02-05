package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// I'm not making two structs for this
type UniqueBody struct {
	maelstrom.MessageBody
	Id uint32 `json:"id,omitempty"`
}

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixMilli())) // this is enough
	node := maelstrom.NewNode()
	node.Handle("generate", func(msg maelstrom.Message) error {
		var body UniqueBody;
		if err := json.Unmarshal(msg.Body,&body); err != nil {
			return err
		}
		body.Type = "generate_ok";
		body.Id = r.Uint32()
		return node.Reply(msg,body)
	});
	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
