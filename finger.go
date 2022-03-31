package main

import (
	"net"
	"net/rpc"
	"net/http"
	"log"
	"fmt"
)

type Node struct{
	SuccMax int
	Port string
	Address string
	Position int
	Fingers [161]string
	Successor []string
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
		n.Successor = []string{fmt.Sprintf("%s:%s", n.Address,n.Port)}
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

func find(id, start string) string {
	found, nextNode := false, start;
	//i = 0;
	for !found /*&& i < maxSteps*/ {
		var output []string
		if err := call(nextNode, "Node.Find", id, &output); err != nil {
			log.Printf("Error finding node")
			return ""
		}
		if len(output[0]) != 0 {
			found = true
		}
		nextNode = output[1]
		//i += 1
	}
	return nextNode
	/*if found {
		return nextNode
	} else {

	}*/
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

func (n *Node) PutAll(kv map[string]string, reply *string) error {
	for key, value := range kv {
		n.Store[key] = value
	}
	*reply = "Values transfered"
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

func (n *Node) GetAll(addr string, reply *[]map[string]string) error {
	ahash := hashstring(addr)
	phash := hashstring(n.Predecessor)
	mop := map[string]string{}
	for key, value := range n.Store {
		khash := hashstring(key)
		if between(phash, khash, ahash, true) {
			mop[key] = value
			delete(n.Store, key)
		}
	}
	*reply = []map[string]string{mop}
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

/*
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
	//n.Predecessor = nil
	//n.Successor = np.find(successor(n)
}*/

func (n *Node) Find(id string, reply *[]string) error {
	b, s := n.find_successor(id)
	var strong string
	if b {
		strong = "true"
	}

	*reply = []string{ strong, s}
	return nil
}
/*
func (n *Node) NewSucc(np string, reply *string) error {
	n.Successor[0] = np
	*reply = "Bitch"
	return nil
}*/

func (n *Node) find_successor(id string) (bool, string){
	idhash := hashstring(id)
	nhash := hashstring(fmt.Sprintf("%s:%s", n.Address, n.Port))
	shash := hashstring(n.Successor[0])
	if between(nhash, idhash, shash, true) {
		return true, n.Successor[0]
	}
	return false, n.Successor[0]
}



func (n *Node) stabilize() {
	firstsucc := n.Successor[0]
	var sp string
	if err := call(firstsucc, "Node.Pred", "", &sp); err != nil {
		if len(n.Successor) <= 1 {
			n.Successor = []string{fmt.Sprintf("%s:%s", n.Address, n.Port)}
		} else {
			n.Successor = n.Successor[1:]
		}
		return
	}
	shash := hashstring(firstsucc)
	sphash := hashstring(sp)
	nhash := hashstring(fmt.Sprintf("%s:%s", n.Address, n.Port))
	if between(nhash, sphash, shash, true){
		n.Successor[0] = sp
	}
	var s []string
	if err := call(n.Successor[0], "Node.Notify", fmt.Sprintf("%s:%s", n.Address, n.Port), &s); err != nil {
		if len(n.Successor) <= 1 {
			n.Successor = []string{fmt.Sprintf("%s:%s", n.Address, n.Port)}
		} else {
			n.Successor = n.Successor[1:]
		}
		return
	}
	if len(s) >= n.SuccMax {
		s = s[:n.SuccMax-1]
	}
	n.Successor = append([]string{n.Successor[0]}, s...)
}

func (n *Node) Pred(none string, reply *string) error {
	*reply = n.Predecessor
	return nil
}

func (n *Node) Notify(np string, reply *[]string) error {
	phash := hashstring(n.Predecessor)
	nphash := hashstring(np)
	nhash := hashstring(fmt.Sprintf("%s:%s", n.Address, n.Port))
	if n.Predecessor == "" || between(phash, nphash, nhash, true) || n.Predecessor == n.Successor[0] {
		n.Predecessor = np
		*reply = n.Successor
	}
	return nil
}
/*
func (n *Node) fix_fingers() {
	next += 1
	if next > m {
		next = 1
	finger[next] = find_successor(
}*/

func (n *Node) check_predecessor() {
	var s string
	if err := call(n.Predecessor, "Node.Failed", "?", &s); err != nil || s == "" {
		n.Predecessor = ""
	}
}

func (n *Node) Failed(none string, reply *string) error {
	*reply = ":)"
	return nil
}

/*
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
