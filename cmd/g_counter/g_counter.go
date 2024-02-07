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

func getCurrentSeqKv(s *maelstrom.KV) (int, error) {
	if value, err := s.ReadInt(context.Background(), KEY); err != nil {
		return 0, err
	} else {
		return value, nil
	}
}

func main() {
	node := maelstrom.NewNode()
	seqkv := maelstrom.NewSeqKV(node)
	node.Handle("read", func(msg maelstrom.Message) error {
		if value, err := seqkv.ReadInt(context.Background(), KEY); err != nil {
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
		value, err := getCurrentSeqKv(seqkv)
		if err == nil {
			// I'm assuming we don't have a counter yet
			// but this read could have failed because we cannot see
			// SeqKV, in which case this would require a retry strategy of some
			// kind
			seqkv.CompareAndSwap(context.Background(), KEY, value, value + body.Delta, true)
		}

		return nil
	})
}
