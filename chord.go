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
	node.SuccMax = 3
	node.Address = getLocalAddress()
	node.Store = make(map[string]string)
	fmt.Printf("Running.\n> ")
	for scan.Scan() {
		text := scan.Text()
		if len(text) != 0 {
			command := strings.Fields(text)
			if funk, ok := commands[command[0]]; ok {
				if !funk(node, command[1:]) {
					fmt.Printf("Command error")
				}
			} else {
				fmt.Printf("Command error")
			}
		} else {
			fmt.Printf("Command error")
		}
		fmt.Printf("\n> ")
	}
}
/*
func in(list []...interface{}, item ...interface{}) bool {
	if len(list) == 0 || reflect.TypeOf(list[0]) != reflect.TypeOf(item) {
		return false
	}
	for _, v := range list {
		if item == v {
			return true
	}
	return false
}*/
