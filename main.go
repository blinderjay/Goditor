package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/blinderjay/Goditor/statik"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/gorilla/websocket"
	"github.com/rakyll/statik/fs"
	"github.com/zserge/webview"
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
	dir                       string      // the path of the execute file
	windowWidth, windowHeight = 1600, 900 // set the window
	statikFS                  http.FileSystem
)

type mdTranser struct {
	sock     *websocket.Conn
	readchn  chan []byte
	writechn chan []byte
}

func init() {
	var err error
	statikFS, err = fs.New()
	if err != nil {
		log.Fatal(err)
	}
	dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	prefixChannel := make(chan string)
	go webServ(prefixChannel)
	prefix := <-prefixChannel
	// create a webview
	wstting := webview.Settings{
		Title: "Goditor",
		// URL:       `data:text/html,` + url.PathEscape(),
		URL:       prefix + "/res/index.html",
		Width:     windowWidth,
		Height:    windowHeight,
		Resizable: true,
	} // the first way to create a window
	w := webview.New(wstting)
	w.Run()
	// err := webview.Open(
	// 	"Goditor",
	// 	prefix+"/res/index.html",
	// 	windowWidth,
	// 	windowHeight,
	// 	false)
	// if err != nil {
	// 	log.Fatal(err)
	// }	// another way to create a window
	log.Printf(prefix)
}
func webServ(prefixChannel chan string) {
	mux := http.NewServeMux()

	/*
		one is used for test front end when all the static file was put in the local;
		while another one embbed all file into go , which is used for final deploy
	*/
	mux.Handle("/res/", http.StripPrefix("/res/", http.FileServer(http.Dir(dir+"/app"))))
	// mux.Handle("/res/", http.StripPrefix("/res/", http.FileServer(statikFS)))
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
				extensions := parser.CommonExtensions |
					parser.MathJax |
					parser.HardLineBreak |
					parser.FencedCode |
					parser.Tables |
					parser.Footnotes |
					parser.Autolink |
					parser.DefinitionLists |
					parser.SuperSubscript |
					parser.Mmark
				parser := parser.NewWithExtensions(extensions)
				md.writechn <- markdown.ToHTML(mdText, parser, nil)
				// md.writechn <- blackfriday.Run(mdText)
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
