// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.



package main

import (
	"flag"
	"github.com/go-basic/uuid"
	"github.com/gorilla/websocket"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)
var menber int = 0
var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options
var clientmap map[string]*websocket.Conn = make(map[string]*websocket.Conn)
var count int = 0
func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	count+=1
	//mt, username, err := c.ReadMessage()
	//log.Println(mt)
	uuid := uuid.New()
	clientmap [uuid] = c
	menber = len(clientmap)
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		for connection := range clientmap{
			if connection ==uuid {
				continue;
			}
			//tgmessage := &pb.Tgmessage{
			//	Username: string(username),
			//	Message: string(message),
			//}
			// ...

			// Write the new address book back to disk.
			//out, err := proto.Marshal(tgmessage)
			//if err != nil {
			//	log.Println("read:", err)
			//	break
			//}
			err = clientmap[connection].WriteMessage(mt, message)
		}
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
	delete(clientmap, uuid)
}

func home(w http.ResponseWriter, r *http.Request) {
	str := strconv.Itoa(menber)
	data := map[string]string{
		"url":    "ws://"+r.Host+"/echo",
		"menber": "在线人数: "+ str,
	}


	homeTemplate.Execute(w, data)
}
//func numb(w http.ResponseWriter, r *http.Request) {
//	homeTemplate.Execute(w, string(menber))
//}
func meassage(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("./tg.js")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	w.Write(content)
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/tg.js", meassage)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))

}


func  ReadHtml() string{
	file, err := os.Open("./index.html")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	return string(content)
}
var homeTemplate = template.Must(template.New("").Parse(ReadHtml()))
