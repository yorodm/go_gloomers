package main

import maelstrom "github.com/jepsen-io/maelstrom/demo/go"


func main(){
	node := maelstrom.NewNode()
	node.Run()
}
