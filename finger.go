package main

import (
	"net"
	"net/rpc"
	"net/http"
	"log"
	"fmt"
)

type Node struct{
	Port string
	Address string
	Position int
	Fingers [161]string
	Successor [1]string
	Predecessor string
	Store map[string]string
	Ring bool
}


func call(address string, method string, request interface{}, reply interface{}) error {
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		return err
	}
	defer client.Close()

	if err = client.Call(method, request, reply); err != nil {
		return err
	}

	return nil
}

func (n *Node) create() error {
	if !n.Ring {
		n.Successor[0] = fmt.Sprintf("%s:%s", n.Address,n.Port)
		n.Predecessor = n.Successor[0]
	}
	rpc.Register(n)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":"+n.Port)
	if e != nil {
		return e
	}
	go http.Serve(l, nil)
	return nil
}

func (n *Node) Ping(p string, reply *string) error {
	fmt.Printf(p)
	*reply = "Pong!"
	return nil
}

func (n *Node) Put(kv []string, reply *string) error {
	n.Store[kv[0]] = kv[1]
	*reply = fmt.Sprintf("%s has been stored under key %s", kv[1], kv[0])
	return nil
}

func (n *Node) Get(key string, reply *string) error {
	if value, ok := n.Store[key]; ok {
		*reply = fmt.Sprintf("%s: %s", key, value)
		return nil
	}
	*reply = "Key not found."
	return nil
}

func (n *Node) Delete(key string, reply *string) error {
	if _, ok := n.Store[key]; ok {
		delete(n.Store, key)
		*reply = fmt.Sprintf("Key/Value: %s, has been deleted", key)
		return nil
	}
	*reply = "Key not found."
	return nil
}


func (n *Node) Join(np string, reply *string) error {
	if n.Predecessor == fmt.Sprintf("%s:%s", n.Address, n.Port) {
		n.Successor[0] = np
	} else {
		var rp string
		if err := call(n.Predecessor, "Node.NewSucc", np, &rp); err != nil {
			log.Printf("Error contacting Predecessor")
			return nil
		}
	}
	*reply = n.Predecessor
	n.Predecessor = np
	return nil
	/*n.Predecessor = nil
	n.Successor = np.find(successor(n)*/
}

func (n *Node) NewSucc(np string, reply *string) error {
	n.Successor[0] = np
	*reply = "Bitch"
	return nil
}
/*
func (n *Node) find_successor(id) {

}
func (n *Node) stabilize() {

}

func (n *Node) notify(np) {

}

func (n *Node) fix_fingers() {

}

func (n *Node) check_predecessor() {

}


func hashString(elt string) *big.Int {
	hasher := sha1.New()
	hasher.Write([]byte(elt))
	return new(big.Int.SetBytes(hasher.Sum(nil))
}*/

func getLocalAddress() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
