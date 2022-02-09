package main

import (
	"fmt"
)

var WsConnections map[string]*WsConnection

type WsConnection struct {
	Uuid        string `json:"Uuid"`
	WsService   string `json:"WsService"`
	StartTime   int64  `json:"StartTime"`
	LastContact int64  `json:"LastContact"`
	IsActive    bool   `json:"IsActive"`
	Conn        string
	ClientConn  string
}

const (
	Register = 0
	Ping     = 1
	Data     = 2
	Error    = 3
	Close    = 9
)

type WsMsg struct {
	Uuid    string `json:"Uuid"`
	Data    []byte `json:"Data"`
	MsgType int    `json:"MsgType"`
}

func GetWsConnections() *map[string]*WsConnection {
	if WsConnections == nil {
		WsConnections = make(map[string]*WsConnection)
	}
	return &WsConnections
}

func main() {
	localwsConnections := GetWsConnections()
	if WsConnections == nil {
		tmpWsConnections := make(map[string]*WsConnection)
		WsConnections = tmpWsConnections
	}
	fmt.Println(localwsConnections)
	localwsConnections = GetWsConnections()
	fmt.Println(localwsConnections)

	c_c2 := &WsConnection{"123", "shell", 1, 1, true, "the agent/c2", ""}
	(*localwsConnections)["123"] = c_c2
	fmt.Println(localwsConnections)
	localwsConnections = GetWsConnections()
	fmt.Println(localwsConnections)

	// fmt.Println(c_c2)

	// registration := WsMsg{"123", []byte("register"), 1}
	// cmd := WsMsg{"123", []byte("whoami"), 2}
	// resp := WsMsg{"123", []byte("root"), 2}
	// _, _ := RegisterWsConnection(registration, "shell", nil, nil)=
	// wsConnections := GetWsConnections()
	// go func() {
	// 	for true {
	// 		if c_c2.Conn != "" {
	// 			fmt.Println(c_c2.Conn)
	// 			c_c2.Conn = ""
	// 		}
	// 		time.Sleep(time.Second * 3)
	// 	}
	// }()
	// for true {
	// 	c_c2.Conn = "HERE"
	// 	time.Sleep(time.Second * 3)
	// }

	// WsConnections[NewConnection.Uuid] = NewConnection
}
