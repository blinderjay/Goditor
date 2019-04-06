package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	blackfriday "gopkg.in/russross/blackfriday.v2"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		// 读取存储空间大小
		ReadBufferSize: 10240,
		// 写入存储空间大小
		WriteBufferSize: 10240,
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	dir                       string     // the path of the execute file
	windowWidth, windowHeight = 800, 600 // set the window
)

type mdTranser struct {
	sock     *websocket.Conn
	readchn  chan []byte
	writechn chan []byte
}

func init() {
	var err error
	dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	prefixChannel := make(chan string)
	webServ(prefixChannel)
	// prefix := <-prefixChannel
	// create a web view
	// err := webview.Open(
	// 	"Goditor",
	// 	prefix+"/res/index.html",
	// 	windowWidth,
	// 	windowHeight,
	// 	false)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Printf(prefix)
}

func webServ(prefixChannel chan string) {
	mux := http.NewServeMux()
	mux.Handle("/res/", http.StripPrefix("/res/", http.FileServer(http.Dir(dir+"/app/res"))))
	mux.HandleFunc("/ws", wsServ)

	listener, err := net.Listen("tcp", "127.0.0.1:8588")
	if err != nil {
		log.Panic(err)
	}
	portAddress := listener.Addr().String()
	listener.Close()
	prefixChannel <- "http://" + portAddress

	server := &http.Server{
		Addr:    portAddress,
		Handler: mux,
	}
	server.ListenAndServe()
}

func wsServ(w http.ResponseWriter, r *http.Request) {
	log.Printf("recerive a websocket connection\n")
	var (
		conn *websocket.Conn
		err  error
	)

	// 完成http应答，在httpheader中放下如下参数
	// Upgrade:websocket
	if conn, err = upgrader.Upgrade(w, r, nil); err != nil {
		log.Println(err) // 获取连接失败直接返回
	}
	md := mdTranser{
		readchn:  make(chan []byte, 409600),
		writechn: make(chan []byte, 409600),
		sock:     conn,
	}

	// 启动转换进程

	go md.getMarkdown()
	go md.sentPreview()

	go func() {
		var mdText []byte
		for {
			select {
			case mdText = <-md.readchn:
				md.writechn <- blackfriday.Run(mdText)
			}
		}
	}()
}

func (md *mdTranser) getMarkdown() {
	defer md.sock.Close()
	for {
		_, mdtext, err := md.sock.ReadMessage()
		if err != nil {
			log.Fatal(err)
			break
		}
		md.readchn <- mdtext
	}
}

func (md *mdTranser) sentPreview() {
	defer md.sock.Close()

	var mdpre []byte
	for {
		select {
		case mdpre = <-md.writechn:
			err := md.sock.WriteMessage(websocket.TextMessage, mdpre)
			if err != nil {
				break
			}
		}
	}
}
