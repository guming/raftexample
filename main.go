package main

import (
	"flag"
	"github.com/coreos/etcd/raft/raftpb"
	"strings"
	"fmt"
)

func main()  {
	cluster:=flag.String("cluster", "http://127.0.0.1:9021", "comma separated cluster peers")
	id := flag.Int("id", 1, "node ID")
	kvport := flag.Int("port", 9121, "key-value server port")

	join := flag.Bool("join", false, "join an existing cluster")
	flag.Parse()
	fmt.Println("kvport:",*kvport)
	var kvs *kvstore

	proposeC := make(chan string)
	defer close(proposeC)
	confChangeC := make(chan raftpb.ConfChange)
	defer close(confChangeC)
	commitC,errorC:=newRaftNode(*id,strings.Split(*cluster,","),*join,proposeC,confChangeC)
	kvs = newKVStore(proposeC, commitC, errorC)
	serveHttpKVAPI(kvs, *kvport, confChangeC, errorC)
}
