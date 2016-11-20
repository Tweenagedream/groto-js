package main

import (
	count "github.com/Tweenagedream/groto-js/protos"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"
	// "io"
	"fmt"
	"net/http"
	"os"
)

func WSServer(ws *websocket.Conn) {
	var buffer []byte
	var counter int32
	websocket.Message.Receive(ws, &buffer)
	helloReq := &count.HelloRequest{}
	proto.Unmarshal(buffer, helloReq)
	counter = helloReq.GetStart()
	fmt.Fprintf(os.Stderr, "Got a new hello: %d\n", counter)
	buffer, _ = proto.Marshal(&count.HelloReply{
		Ack: *proto.Bool(true),
	})
	websocket.Message.Send(ws, buffer)

}

func main() {
	http.Handle("/ws/", websocket.Handler(WSServer))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.Handle("/static/", http.FileServer(http.Dir("protos")))

	http.ListenAndServe(":8000", nil)
}
