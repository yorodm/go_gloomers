package main

import (
	"context"
	"encoding/json"
	"fmt"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type ReadResponse struct {
	maelstrom.MessageBody
	Value int `json:"value,omitempty"`
}

type AddBody struct {
	maelstrom.MessageBody
	Delta int `json:"delta,omitempty"`
}

const KEY = "my_g_counter"

func main() {
	node := maelstrom.NewNode()
	seqkv := maelstrom.NewSeqKV(node)
	state := 0
	node.Handle("read", func(msg maelstrom.Message) error {
		if value, err := seqkv.ReadInt(context.TODO(), KEY); err != nil {
			return err
		} else {
			response := ReadResponse{}
			response.Type = "read_ok"
			response.Value = value
			return node.Reply(msg, response)

		}
	})
	node.Handle("add", func(msg maelstrom.Message) error {
		var body AddBody
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		if err := seqkv.CompareAndSwap(context.TODO(), KEY, state, state+body.Delta, true); err != nil {
			// READ and update state
		}
		return nil
	})
}
