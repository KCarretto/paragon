package websockets

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var PingTime = time.Second * 10
var PingGrace = time.Second * 5
var WsConnections map[string]*WsConnection
var ValidWsServices = [2]string{"shell", "socks"}

type WsConnection struct {
	Uuid        string `json:"Uuid"`
	WsService   string `json:"WsService"`
	StartTime   int64  `json:"StartTime"`
	LastContact int64  `json:"LastContact"`
	IsActive    bool   `json:"IsActive"`
	Conn        *websocket.Conn
	ClientConn  *websocket.Conn
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

func GetWsConnection(uuid string) (*WsConnection, error) {
	wsConnections := GetWsConnections()
	// Get connection from map.
	if val, ok := (*wsConnections)[uuid]; ok {
		wsConnection := val
		// Update active status
		wsConnection.IsActive = (wsConnection.LastContact+int64(PingTime)+int64(PingGrace) < time.Now().Unix())
		return wsConnection, nil
	}
	return nil, errors.New("could not find WsConnection for: ")

}

func ConnectionPing(uuid string) error {
	connection, err := GetWsConnection(uuid)
	if err != nil {
		return err
	}
	if connection != (nil) {
		connection.LastContact = time.Now().Unix()
		return nil
	}
	return errors.New("connection not found may need to register")
}

func RegisterWsConnection(wsMsg *WsMsg, WsService string, Conn *websocket.Conn, logger *zap.Logger) (*WsConnection, error) {
	wsConnection := GetWsConnections()
	isValidService := false
	for _, value := range ValidWsServices {
		if value == WsService {
			isValidService = true
			break
		}
	}
	if isValidService == false {
		return nil, errors.New(fmt.Sprintf("Not a valid WsService: %s\nPlease use one of: %v", WsService, ValidWsServices))
	}
	logger.Debug("is a valid service")
	NewConnection := WsConnection{
		Uuid:        wsMsg.Uuid,
		WsService:   WsService,
		Conn:        Conn,
		ClientConn:  nil,
		StartTime:   time.Now().Unix(),
		LastContact: time.Now().Unix(),
		IsActive:    true,
	}
	logger.Debug("created conn object")
	(*wsConnection)[NewConnection.Uuid] = &NewConnection
	logger.Debug("stored conn object")
	return &NewConnection, nil
}

func (wm *WsMsg) ToJson() ([]byte, error) {
	WsMsgString, err := json.Marshal(wm)
	if err != nil {
		return []byte("Failed to encode WsMsg into JSON string"), err
	}
	return []byte(WsMsgString), nil
}

func ToStruct(WsMsgString []byte) (*WsMsg, error) {
	WsMsgStruct := &WsMsg{}
	err := json.Unmarshal(WsMsgString, WsMsgStruct)
	if err != nil {
		return nil, err
	}
	return WsMsgStruct, nil
}
