package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"html/template"
	"net/http"
	"os"
)

type position struct {
	status int64
	x int64
	y int64
}

type user_conn struct {
	ws *websocket.Conn
	status int64
}

var user_list []*user_conn

func conn(ws *websocket.Conn){
	//添加到连接队列里面
	res := addqueue(ws)
	if res < 0 {
		websocket.Message.Send(ws,"已有两个人在线!!!")
		return
	}
			for {
				var reply string
				if err := websocket.Message.Receive(ws, &reply); err != nil {
					ws.Close()
					delqueue(ws)
					return
				}
				for _, b := range user_list {
					if b.ws != ws {
						if er := websocket.Message.Send(b.ws, reply); er != nil {
							b.ws.Close()
						}
					}
				}
			}
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)

}

func addqueue(ws *websocket.Conn) int {
	for a,b := range user_list{
		if  b.ws == ws{
			return a
		}
	}
	current_link := &user_conn{ws,1}
	l := len(user_list)
	if l >= 2 {
		return -1
	}
	user_list = append(user_list, current_link)
	return l
}

func delqueue(ws *websocket.Conn)  {
	for a,b := range user_list{
		if  b.ws == ws{
			user_list = append(user_list[0:a],user_list[a+1:]...)
		}
	}
	fmt.Println(user_list)
}


func main() {
	http.Handle("/conn", websocket.Handler(conn))
	http.HandleFunc("/", index)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

