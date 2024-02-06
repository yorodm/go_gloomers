package main

import (
	"encoding/json"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type ReadBody struct {
	maelstrom.MessageBody
	Messages []int `json:"messages,omitempty"`
}

type BroadcastBody struct {
	maelstrom.MessageBody
	Value int `json:"messages,omitempty"`
}

type TopologyBody struct {
	maelstrom.MessageBody
	Topology map[string][]string `json:"topology,omitempty"`
}

type GossipBody struct {
	maelstrom.MessageBody
	Values []int `json:"values,omitempty"`
}

type void struct {}

type State struct {
	mut          sync.RWMutex
	messages map[int]void
}

//NewState: create a new State
func NewState() *State {
	return &State{messages: make(map[int]void)}
}

func (s *State) Messages() []int {
	s.mut.RLock()
	defer s.mut.RUnlock()
	data := make([]int, len(s.messages)) // speed up
	for key, _ := range s.messages {
		data = append(data, key)
	}
	return data
}

func (s *State)Add(v int) {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.messages[v] = void{}
}

func (s *State)AddAll(v []int) {
	s.mut.Lock()
	defer s.mut.Unlock()
	for _, e := range v {
		s.messages[e] = void{}
	}
}

func deserialize(raw json.RawMessage, v any) error {
	if err := json.Unmarshal(raw, v); err != nil {
		return err
	}
	return nil
}

func main() {
	var neighbours []string
	state := NewState()
	node := maelstrom.NewNode()
	node.Handle("read", func(msg maelstrom.Message) error {
		var body ReadBody
		if err := deserialize(msg.Body, &body); err != nil {
			return err
		}
		body.Type = "read_ok"
		body.Messages = state.Messages()

		return node.Reply(msg, body)
	})
	node.Handle("broadcast", func(msg maelstrom.Message) error {
		var body BroadcastBody
		body.Type = "broadcast_ok"
		go func() {
			b := GossipBody {}
			b.Type = "gossip"
			b.Values = state.Messages()
			for _, n := range neighbours {
				node.Send(n, b)
			}
		}()
		return node.Reply(msg, &body)
	})
	node.Handle("topology", func(msg maelstrom.Message) error {
		var body TopologyBody
		if err := deserialize(msg.Body, &body); err != nil {
			return err
		}
		 neighbours = body.Topology[node.ID()]
		body.Type = "topology_ok"
		return node.Reply(msg, body)
	})
	node.Handle("gossip", func(msg maelstrom.Message) error {
		var body GossipBody
		if err := deserialize(msg.Body, &Body)
	})
}
