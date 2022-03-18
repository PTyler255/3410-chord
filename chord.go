package main

import (
	"fmt"
	//"log"
	"bufio"
	"os"
	"strings"
	//"crypto/sha1"
	//"math/big"
	//"net/rpc"
)

//map[string]func(*Node, []string)bool


func main() {
	commands := make(map[string]func(*Node, []string)bool)
	addCommands(commands)
	scan := bufio.NewScanner(os.Stdin)
	node := new(Node)
	node.Port = "3410"
	node.Address = getLocalAddress()
	node.Store = make(map[string]string)
	fmt.Printf("Running.\n> ")
	for scan.Scan() {
		command := strings.Fields(scan.Text())
		if funk, ok := commands[command[0]]; ok {
			if !funk(node, command[1:]) {
				fmt.Printf("Command error")
			}
		} else {
			fmt.Printf("Command error")
		}
		fmt.Printf("\n> ")
	}
}


