package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

var WsConnections map[string]WsConnection

type WsConnection struct {
	Uuid        string `json:"Uuid"`
	WsService   string `json:"WsService"`
	StartTime   int64  `json:"StartTime"`
	LastContact int64  `json:"LastContact"`
	IsActive    bool   `json:"IsActive"`
	Conn        *websocket.Conn
}

func main() {
	c_c2 := WsConnection{}
	if c_c2.Uuid == nil {
		fmt.Println("Not set")
	}
	WsConnections[NewConnection.Uuid] = NewConnection
}
