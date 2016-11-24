package main

import (
	count "github.com/Tweenagedream/groto-js/protos"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"
	// "io"
	"fmt"
	"net/http"
	"os"
	"sync"
)

var conns map[int32]*websocket.Conn
var id, gCount int32
var idLock sync.Mutex
var countLock sync.Mutex

func getId() int32 {
	idLock.Lock()
	defer idLock.Unlock()
	ret := id
	id++
	return ret
}

func broadcastToConns(data []byte, err error) error {
	for conn := range conns {
		err = websocket.Message.Send(conns[conn], data)
	}
	return err
}

func addToCount(addition int32) error {
	countLock.Lock()
	defer countLock.Unlock()
	gCount += addition
	fmt.Fprintf(os.Stderr, "New gCount value: %v\n", gCount)
	creply := &count.CountReply{
		Count: gCount,
	}
	err := broadcastToConns(proto.Marshal(creply))
	return err
}

func deleteAndLog(connId int32) {
	fmt.Fprintf(os.Stderr, "Deleting connection: %v\n", connId)
	delete(conns, connId)
}

func WSServer(ws *websocket.Conn) {
	var buffer []byte
	var counter int32
	connId := getId()
	conns[connId] = ws
	defer deleteAndLog(connId)
	fmt.Fprintf(os.Stderr, "New connection: %v\n", connId)
	websocket.Message.Receive(ws, &buffer)
	helloReq := &count.HelloRequest{}
	proto.Unmarshal(buffer, helloReq)
	counter = helloReq.GetStart()
	fmt.Fprintf(os.Stderr, "Got a new hello: %d\n", counter)
	buffer, _ = proto.Marshal(&count.HelloReply{
		Ack: *proto.Bool(true),
	})
	websocket.Message.Send(ws, buffer)
	countRequest := &count.CountRequest{}
	ok := true
	var err error
	for ok {
		websocket.Message.Receive(ws, &buffer)
		err = proto.Unmarshal(buffer, countRequest)
		fmt.Fprintf(os.Stderr, "Value of buffer: <%v>\n", buffer)
		if err != nil {
			ok = false
			break
		}
		if countRequest.GetCount() > 0 {
			err = addToCount(countRequest.GetCount())
			if err != nil {
				ok = false
				break
			}
		}
	}

}

func main() {
	conns = make(map[int32]*websocket.Conn)
	http.Handle("/ws/", websocket.Handler(WSServer))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.Handle("/static/", http.FileServer(http.Dir("protos")))

	http.ListenAndServe(":8000", nil)
	fmt.Fprintf(os.Stderr, "Listening on port 8000...")
}
