package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/gorilla/websocket"
)

var (
	uupgrader = websocket.Upgrader{
		// 读取存储空间大小
		ReadBufferSize: 10240,
		// 写入存储空间大小
		WriteBufferSize: 10240,
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ddir string
)

func TestWeb(t *testing.T) {
	var err error
	ddir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("testing Web")
	mux := http.NewServeMux()
	mux.Handle("/res/", http.StripPrefix("/res/", http.FileServer(http.Dir(ddir+"/app/"))))
	mux.HandleFunc("/ws", wsServ)
	mux.HandleFunc("/", toindex)
	listener, err := net.Listen("tcp", "127.0.0.1:8588")
	if err != nil {
		log.Panic(err)
	}
	portAddress := listener.Addr().String()
	listener.Close()
	server := &http.Server{
		Addr:    portAddress,
		Handler: mux,
	}
	server.ListenAndServe()
}

func toindex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/res/index.html", 302)
}
