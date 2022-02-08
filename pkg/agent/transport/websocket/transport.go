package websocket

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/kcarretto/paragon/pkg/agent/transport"
	"go.uber.org/zap"
)

type ServerTransport struct {
	transport.Transport
	Log        *zap.Logger
	Handler    func(w http.ResponseWriter, req *http.Request, t *ServerTransport)
	TeamServer string
}

func (t *ServerTransport) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Setup request logging
	// w.Write([]byte("The time is: " + "Now"))
	t.Log.Debug("test")
	t.Handler(w, req, t)

}

func HandleGiveShell(w http.ResponseWriter, req *http.Request, t *ServerTransport) {
	u := url.URL{Scheme: "ws", Host: t.TeamServer, Path: "websocketgiveshell"}
	// w.Write([]byte("The time is: " + u.String()))

	var upgrader = websocket.Upgrader{} // use default options
	c_agent, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		t.Log.Debug(fmt.Sprintf("upgrade: %v\n", err))
		return
	}

	t.Log.Debug(fmt.Sprintf("Trying to connect to:%s\n", u.String()))
	c_teamserver, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Log.Debug(fmt.Sprintf("dial: %v", err))
		return
	}
	defer c_teamserver.Close()
	t.Log.Debug("Called out")

	defer c_agent.Close()
	// handle agent messages
	go func() {
		for {
			mt, message, err := c_agent.ReadMessage()
			if err != nil {
				t.Log.Debug(fmt.Sprintf("read: %v\n", err))
				break
			}

			t.Log.Info(string(message))
			log.Printf("recv: %s", message)
			err = c_teamserver.WriteMessage(mt, message)
			if err != nil {
				t.Log.Debug(fmt.Sprintf("write: %v\n", err))
				break
			}
		}
	}()

	// // handle client messages
	for {
		mt, message, err := c_teamserver.ReadMessage()
		if err != nil {
			t.Log.Debug(fmt.Sprintf("read: %v\n", err))
			break
		}

		t.Log.Info(string(message))
		log.Printf("recv: %s", message)
		err = c_agent.WriteMessage(mt, message)
		if err != nil {
			t.Log.Debug(fmt.Sprintf("write: %v\n", err))
			break
		}
	}

}
