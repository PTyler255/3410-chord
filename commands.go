package main

import (
	"os"
	"fmt"
	"strconv"
	"log"
)

//map[string]func(*Node, []string)bool

func addCommands(commands map[string]func(*Node, []string)bool) {
	commands["help"] = doHelp
	commands["port"] = doPort
	commands["create"] = doCreate
	commands["join"] = doJoin
	commands["ping"] = doPing
	commands["getip"] = doGetIP
	commands["put"] = doPut
	commands["putrandom"] = doPutRandom
	commands["get"] = doGet
	commands["delete"] = doDelete
	commands["dump"] = doDump
	commands["dumpkey"] = doDumpKey
	commands["dumpaddr"] = doDumpAddr
	commands["dumpall"] = doDumpAll
	commands["quit"] = doQuit
	commands["q"] = doQuit
}

func doHelp(n *Node, none []string) bool {
	if len(none) != 0 {
		return false
	}
	fmt.Printf("- help\n- port <n>\n- create\n- ping\n- getip\n- join <address>\n- put <key> <value>\n- putrandom <n>\n- get <key>\n- delete <key>\n- dump\n- dumpkey <key>\n- dumpaddr <address>\n- dumpall\n- quit")
	return true
}

func doPort(n *Node, pn []string) bool {
	if len(pn) != 1 {
		return false
	}
	if _, err := strconv.Atoi(pn[0]); err != nil {
		return false
	}
	n.Port = pn[0]
	fmt.Printf("Port number changed to: %s", n.Port)
	return true
}

func doCreate(n *Node, none []string) bool {
	if len(none) != 0 || n.Ring {
		return false
	}
	if err := n.create(); err != nil {
		log.Printf("Creation error: %v", err)
		return false
	}
	n.Ring = true
	return true
}

func doJoin(n *Node, addr []string) bool {
	if len(addr) != 1 || n.Ring {
		return false
	}
	var s string
	if err := call(addr[0], "Node.Join", fmt.Sprintf("%s:%s", n.Address,n.Port), &s); err != nil {
		log.Printf("Join error: %v", err)
		return false
	}
	if s != "" {
		n.Successor[0] = addr[0]
		n.Predecessor = s
		n.Ring = true
		if err := n.create(); err != nil {
			log.Printf("Creation error: %v", err)
			return false
		}
		return true
	}
	return false
}

func doGetIP(n *Node, none []string) bool {
	if len(none) != 0 {
		return false
	}
	fmt.Printf(getLocalAddress())
	return true
}

func doPing(n *Node, addr []string) bool {
	if len(addr) != 1 {
		return false
	}
	var s string
	if err := call(addr[0], "Node.Ping", "ping", &s); err != nil {
		log.Printf("Error with call: %v", err)
		return false
	}
	fmt.Printf(s)
	return true
}

//-------------------------------------------------------

func doPut(n *Node, args []string) bool {
	if len(args) != 3 || !n.Ring {
		return false
	}
	var s string
	if err := call(args[2], "Node.Put", args[:2], &s); err != nil {
		log.Printf("Error with call: %v", err)
		return false
	}
	fmt.Printf(s)
	return true
}

func doPutRandom(n *Node, addr []string) bool {
	return false

}

func doGet(n *Node, args []string) bool {
	if len(args) != 2 || !n.Ring {
		return false
	}
	var s string
	if err := call(args[1], "Node.Get", args[0], &s); err != nil {
		log.Printf("Error with call: %v", err)
		return false
	}
	fmt.Printf(s)
	return true
}

func doDelete(n *Node, args []string) bool {
	if len(args) != 2 || !n.Ring {
		return false
	}
	var s string
	if err := call(args[1], "Node.Delete", args[0], &s); err != nil {
		log.Printf("Error with call: %v", err)
		return false
	}
	fmt.Printf(s)
	return true
}

//-------------------------------------------------------

func doDump(n *Node, none []string) bool {
	if len(none) != 0 {
		return false
	}
	fmt.Printf("Dumping Store.")
	for k, v := range n.Store {
		fmt.Printf("\nKey: %s\tValue: %s", k,v)
	}
	fmt.Printf("\nAll Successors.")
	for i, v := range n.Successor {
		fmt.Printf("\n%d: %s", i, v)
	}
	return true

}

func doDumpKey(n *Node, addr []string) bool {
	return false

}

func doDumpAddr(n *Node, addr []string) bool {
	return false

}

func doDumpAll(n *Node, none []string) bool {
	if len(none) != 0 {
		return false
	}
	return false

}

//-------------------------------------------------------

func doQuit(n *Node, none []string) bool {
	if len(none) != 0 {
		return false
	}
	os.Exit(0)
	return true
}
