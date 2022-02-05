package websockets

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

const (
	Registration int = 0
	Command          = 1
	Response         = 2
)

const (
	Agent   int = 0
	Client      = 1
	Control     = 2
)

type WsConn struct {
	Conn   *websocket.Conn
	Send   chan []WsMsg //May need to back to being byte.
	Uuid   string
	Active bool
}

type WsMsg struct {
	Uuid    string
	Data    string
	SrcType int
	MsgType int
}

func (wm WsMsg) ToJson() (string, error) {
	out_string, err := json.Marshal(WsMsg{Uuid: string(wm.Uuid), Data: string(wm.Data), SrcType: wm.SrcType, MsgType: wm.MsgType})
	if err != nil {
		log.Println("ToString:", err)
		return "", err
	}
	return string(out_string), err
}

func (wm WsMsg) ToString() (string, error) {
	res := fmt.Sprintf("WsMsg {"+
		"\n	wm.Uuid:		%s"+
		"\n	wm.Data:		%s"+
		"\n	wm.SrcType:		%d"+
		"\n	wm.MsgType:		%d"+
		"\n}", wm.Uuid, wm.Data, wm.SrcType, wm.MsgType)
	return string(res), nil
}

func WsMsgFromJson(jsonString string, wm *WsMsg) error {
	err := json.Unmarshal([]byte(jsonString), wm)
	if err != nil {
		log.Println("JSON conversion")
		log.Println(err)
	}

	return err
}

func TestToJson() {
	validMatch := "{\"Uuid\":\"test1\",\"Data\":\"Test2\",\"SrcType\":0,\"MsgType\":0}"
	wsMsg := WsMsg{Uuid: "test1", Data: "Test2", SrcType: Agent}
	jsonRes, err := wsMsg.ToJson()
	if err != nil {
		fmt.Println("TestToJson 			err - ToJson()")
	}
	if string(jsonRes) == validMatch {
		fmt.Println("TestToJson 			passed")
	} else {
		fmt.Println("TestToJson 			err - content match")
		fmt.Println(string(jsonRes))
		fmt.Println(validMatch)
	}
}

func TestWsMsgFromJson() {
	wsMsg := &WsMsg{}
	err := WsMsgFromJson("{\"Uuid\":\"test1\",\"Data\":\"Test2\",\"SrcType\":1,\"MsgType\":0}", wsMsg)
	if err != nil {
		fmt.Println("TestWsMsgFromJson		err - ToJson()")
	}
	if wsMsg.Uuid == "test1" && wsMsg.Data == "Test2" && wsMsg.SrcType == 1 && wsMsg.MsgType == 0 {
		fmt.Println("TestWsMsgFromJson		passed")
	} else {
		fmt.Println("TestWsMsgFromJson		err - content match")
		fmt.Println(wsMsg.ToString())
	}
}

func Test() {
	TestToJson()
	TestWsMsgFromJson()
}
