package main

import (
    "encoding/json"
    "log"

    maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// Follow the lead from the library, simpler than dealing with
// all kinds of maps and reclection and whatnot
type EchoBody struct {
	maelstrom.MessageBody
	Echo string `json:"echo,omitempty"`
}

func main() {
	node := maelstrom.NewNode()
	node.Handle("echo", func(msg maelstrom.Message) error {
		var body EchoBody
		if err := json.Unmarshal(msg.Body,&body); err != nil {
			return err
		}
		body.Type = "echo_ok";
		return node.Reply(msg, body)
	})
	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
