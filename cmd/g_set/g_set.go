package main

import (
	"encoding/json"
	"log"
	"slices"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type ReadResponseBody struct {
	maelstrom.MessageBody
	Value []int `json:"messages"`
}

type AddBody struct {
	maelstrom.MessageBody
	Element int `json:"element,omitempty"`
}

type GossipBody struct {
	maelstrom.MessageBody
	Values []int `json:"values"`
}

type void struct{}

type State struct {
	mut      sync.RWMutex
	messages map[int]void
}

// NewState: create a new State
func NewState() *State {
	return &State{messages: make(map[int]void)}
}

func (s *State) Messages() []int {
	s.mut.RLock()
	defer s.mut.RUnlock()
	data := make([]int, 0, len(s.messages)) // speed up
	for key := range s.messages {
		data = append(data, key)
	}
	slices.Sort(data)
	return data
}

func (s *State) Add(v int) {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.messages[v] = void{}
}

func (s *State) AddAll(v []int) {
	s.mut.Lock()
	defer s.mut.Unlock()
	for _,e := range v {
		s.messages[e] = void{}
	}
}

func deserialize(raw json.RawMessage, v any) error {
	if err := json.Unmarshal(raw, v); err != nil {
		return err
	}
	return nil
}

func gossip(node *maelstrom.Node, value []int, neighbours []string) {
	b := GossipBody{}
	b.Type = "gossip"
	b.Values = value
	for _, n := range neighbours {
		node.Send(n, b)
	}
}

func main() {
	var neighbours []string
	state := NewState()
	node := maelstrom.NewNode()
	node.Handle("read", func(msg maelstrom.Message) error {
		body := ReadResponseBody{
			Value: state.Messages(),
		}
		body.Type = "read_ok"
		return node.Reply(msg, &body)
	})
	node.Handle("add", func(msg maelstrom.Message) error {
		var body BroadcastBody
		if err := deserialize(msg.Body, &body); err != nil {
			return err
		}
		state.Add(body.Message)
		log.Printf("Vecinos %s", neighbours)
		go gossip(node, state.Messages(), neighbours)
		reply := map[string]string{
			"type": "broadcast_ok",
		}
		return node.Reply(msg, &reply)
	})
	node.Handle("gossip", func(msg maelstrom.Message) error {
		var body GossipBody
		if err := deserialize(msg.Body, &body); err != nil {
			log.Fatal(err)
		}
		state.AddAll(body.Values)
		rest := make([]string,0)
		for _, n := range neighbours {
			if n != msg.Src {
				rest = append(rest, n)
			}
		}
		go gossip(node,state.Messages(),rest)
		return nil
	})
	node.Run()
}
