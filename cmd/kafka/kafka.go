package main

import maelstrom "github.com/jepsen-io/maelstrom/demo/go"

type SendBody struct {
	Key string `json:"key,omitempty"`
	Msg uint64 `json:"msg,omitempty"`
}


func main() {
	node := maelstrom.NewNode()

	node.Run()
}
