package teamserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/graphql"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"github.com/kcarretto/paragon/pkg/auth"
	"github.com/kcarretto/paragon/pkg/auth/oauth"
	"github.com/kcarretto/paragon/pkg/cdn"
	"github.com/kcarretto/paragon/pkg/event"
	"github.com/kcarretto/paragon/pkg/service"
	commonwebsockets "github.com/kcarretto/paragon/pkg/websockets"
	"github.com/kcarretto/paragon/www"
)

// Service provides HTTP handlers to compose the CDN, GraphQL, and WWW services.
type Service struct {
	Log    *zap.Logger
	Graph  *ent.Client
	Events event.Publisher
	OAuth  *oauth2.Config
	Auth   service.Authenticator
}

// HandleStatus returns JSON status: OK if the teamserver is running without error.
func (svc *Service) HandleStatus(w http.ResponseWriter, r *http.Request) error {
	type Response struct {
		Available       bool `json:"available"`
		UserID          int  `json:"userid,omitempty"`
		IsAuthenticated bool `json:"is_authenticated"`
		IsActivated     bool `json:"is_activated"`
		IsAdmin         bool `json:"is_admin"`
	}
	dst := json.NewEncoder(w)

	resp := Response{
		Available: true,
	}
	if user := auth.GetUser(r.Context()); user != nil {
		resp.UserID = user.ID
		resp.IsAuthenticated = true
		resp.IsActivated = user.IsActivated
		resp.IsAdmin = user.IsAdmin
	}

	if err := dst.Encode(resp); err != nil {
		return fmt.Errorf("failed to write response data: %w", err)
	}
	w.Header().Set("Content-Type", "application/json")
	return nil
}

// HTTP registers http handlers for the Teamserver.
func (svc *Service) HTTP(router *http.ServeMux) {
	oauthSVC := &oauth.Service{
		Log:    svc.Log.Named("auth"),
		Graph:  svc.Graph,
		Config: svc.OAuth,
	}
	graphqlSVC := &graphql.Service{
		Log:    svc.Log.Named("graphql"),
		Graph:  svc.Graph,
		Events: svc.Events,
		Auth:   svc.Auth,
	}
	cdnSVC := &cdn.Service{
		Log:   svc.Log.Named("cdn"),
		Graph: svc.Graph,
		Auth:  svc.Auth,
	}
	wwwSVC := &www.Service{
		Log:  svc.Log.Named("www"),
		Auth: svc.Auth,
	}
	status := &service.Endpoint{
		Log:           svc.Log.Named("status"),
		Authenticator: svc.Auth,
		Handler:       service.HandlerFn(svc.HandleStatus),
	}
	initSVC := &InitService{
		Log:   svc.Log.Named("init"),
		Graph: svc.Graph,
		Auth:  svc.Auth,
	}
	wsShell := &service.Endpoint{
		Log:           svc.Log.Named("wsShell"),
		Authenticator: svc.Auth,
		Handler:       service.HandlerFn(svc.HandleGiveShell),
	}
	wsShellConnect := &service.Endpoint{
		Log:           svc.Log.Named("wsShellConnect"),
		Authenticator: svc.Auth,
		Handler:       service.HandlerFn(svc.HandleShellConnect),
	}

	router.Handle("/websocketgiveshell", wsShell)
	router.Handle("/websocketconnectshell", wsShellConnect)

	graphqlSVC.HTTP(router)
	cdnSVC.HTTP(router)
	wwwSVC.HTTP(router)

	router.Handle("/status", status)
	initSVC.HTTP(router)
	oauthSVC.HTTP(router)
}

func (svc *Service) HandleShellConnect(w http.ResponseWriter, r *http.Request) error {
	svc.Log.Debug("trying to connect shell")
	var upgrader = websocket.Upgrader{} // use default options

	// recieve forwarded message from c2
	c_c2 := &commonwebsockets.WsConnection{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c_client, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		svc.Log.Debug(fmt.Sprintf("upgrade: %v\n", err))
		return err
	}
	defer c_client.Close()
	wsConnections := commonwebsockets.GetWsConnections()
	//check the current wsConnections size excluding the new one being made.
	svc.Log.Debug(fmt.Sprintf("wsConnections length: %d", len(*wsConnections)))

	for {
		// Get command from client
		_, message, err := c_client.ReadMessage()
		if err != nil {
			svc.Log.Debug(fmt.Sprintf("read: %v\n", err))
			break
		}
		svc.Log.Info(fmt.Sprintf("recieved data: %s", message))

		// Handle websocket
		wsMsg, err := commonwebsockets.ToStruct(message)
		if err != nil {
			svc.Log.Debug(fmt.Sprintf("JSON Unmarshal: %v\n", err))
			continue
		}

		// Handle the message
		switch wsMsg.MsgType {
		// case commonwebsockets.Register:
		// 	svc.Log.Debug(fmt.Sprintf("Registering client %s", wsMsg.Uuid))
		// 	// Find active connection for the agent
		// 	c_c2, err = commonwebsockets.GetWsConnection(wsMsg.Uuid)
		// 	if err != nil {
		// 		svc.Log.Debug(fmt.Sprintf("Couldn't find client connection  %s", err))
		// 		continue
		// 	}
		// 	if c_c2 == (commonwebsockets.WsConnection{}) {
		// 		svc.Log.Debug(fmt.Sprintf("Couldn't find client connection  %s", wsMsg.Uuid))
		// 		continue
		// 	}
		// 	c_c2.ClientConn = c_client
		// 	// _, err := commonwebsockets.RegisterWsConnection(wsMsg, "shell", c_c2, svc.Log)
		// 	// if err != nil {
		// 	// 	svc.Log.Debug(fmt.Sprintf("Failed to register agent: %v\n", err))
		// 	// 	continue
		// 	// }
		// 	svc.Log.Debug(fmt.Sprintf("Updating client ping  %s", wsMsg.Uuid))
		// 	commonwebsockets.ConnectionPing(wsMsg.Uuid)
		case commonwebsockets.Data:
			svc.Log.Debug(fmt.Sprintf("Recieved request command %s", wsMsg.Data))
			c_c2, err = commonwebsockets.GetWsConnection(wsMsg.Uuid)
			if err != nil {
				svc.Log.Debug(fmt.Sprintf("Error getting ws conn for agent %v\n", err))
				continue
			}
			c_c2.ClientConn = c_client
			// fmt.Println("c_client: %v", c_client)
			// fmt.Println("c_c2.ClientConn: %v", c_c2.ClientConn)
			c_c2, err = commonwebsockets.GetWsConnection(wsMsg.Uuid)
			// fmt.Println("refreshed c_c2.ClientConn: %v", c_c2.ClientConn)

			if c_c2 != (nil) {
				wsMsgJSON, err := wsMsg.ToJson()
				if err != nil {
					svc.Log.Debug(fmt.Sprintf("Error converting wsMsg to JSON %v\n", err))
					continue
				}
				c_c2.Conn.WriteMessage(websocket.TextMessage, wsMsgJSON)
			} else {
				svc.Log.Debug(fmt.Sprintf("c_c2 conn is not set yet please register first %s", wsMsg.Data))
				continue
			}
		case commonwebsockets.Ping:
			svc.Log.Debug(fmt.Sprintf("Recieved ping %s", string(wsMsg.Uuid)))
			err := commonwebsockets.ConnectionPing(wsMsg.Uuid)
			if err != nil {
				svc.Log.Debug(fmt.Sprintf("Could not find conn for %s", string(wsMsg.Uuid)))
				continue
			}
		case commonwebsockets.Error:
			svc.Log.Debug(fmt.Sprintf("Recieved error %s", wsMsg.Data))
		case commonwebsockets.Close:
			svc.Log.Debug(fmt.Sprintf("Recieved close %s", wsMsg.Uuid))

		}
	}
	return nil
}

func (svc *Service) HandleGiveShell(w http.ResponseWriter, r *http.Request) error {
	svc.Log.Debug("trying to give shell")
	var upgrader = websocket.Upgrader{} // use default options

	// recieve forwarded message from c2
	c_c2, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		svc.Log.Debug(fmt.Sprintf("upgrade: %v\n", err))
		return err
	}
	defer c_c2.Close()
	wsConnections := commonwebsockets.GetWsConnections()
	//check the current wsConnections size excluding the new one being made.
	svc.Log.Debug(fmt.Sprintf("wsConnections length: %d", len(*wsConnections)))

	for {
		// Read message from c2
		_, message, err := c_c2.ReadMessage()
		if err != nil {
			svc.Log.Debug(fmt.Sprintf("read: %v\n", err))
			break
		}
		svc.Log.Info(fmt.Sprintf("recieved data: %s", message))

		// Handle websocket
		wsMsg, err := commonwebsockets.ToStruct(message)
		if err != nil {
			svc.Log.Debug(fmt.Sprintf("JSON Unmarshal: %v\n", err))
			continue
		}

		// Handle the message
		switch wsMsg.MsgType {
		case commonwebsockets.Register:
			svc.Log.Debug(fmt.Sprintf("Registering client %s", wsMsg.Uuid))
			_, err := commonwebsockets.RegisterWsConnection(wsMsg, "shell", c_c2, svc.Log)
			if err != nil {
				svc.Log.Debug(fmt.Sprintf("Failed to register agent: %v\n", err))
				continue
			}
			svc.Log.Debug(fmt.Sprintf("Updating client ping  %s", wsMsg.Uuid))
			commonwebsockets.ConnectionPing(wsMsg.Uuid)
			// // DEBUG
			// svc.Log.Debug(fmt.Sprintf("Building agent response (debug)  %s", wsMsg.Uuid))
			// respWsMsg := commonwebsockets.WsMsg{Uuid: wsMsg.Uuid, Data: []byte("whoami"), MsgType: commonwebsockets.Data}
			// respWsMsgJSON, err := respWsMsg.ToJson()
			// if err != nil {
			// 	svc.Log.Debug(fmt.Sprintf("Failed to marshal JSON: %v\n", err))
			// 	continue
			// }
			// c_c2.WriteMessage(websocket.TextMessage, respWsMsgJSON)
			// // END DEBUG

		case commonwebsockets.Data:
			// ERROR APPEARS TO BE IN HERE
			// runtime error: invalid memory address or nil pointer dereference
			// wsConnections := commonwebsockets.GetWsConnections()
			wsConn, err := commonwebsockets.GetWsConnection(wsMsg.Uuid)
			if err != nil {
				svc.Log.Debug(fmt.Sprintf("Could not find conn for %s", string(wsMsg.Uuid)))
				continue
			}
			svc.Log.Debug(fmt.Sprintf("Was able to get agent websocket conn"))

			wsMsgJSON, err := wsMsg.ToJson()
			svc.Log.Debug(fmt.Sprintf("ToJson"))
			svc.Log.Debug(fmt.Sprintf("STARTING BUG ANALYSIS"))
			// looks like clientConn is nil.
			// Error check that it's not nil.

			if wsConn.ClientConn == nil {
				svc.Log.Debug(fmt.Sprintf("Could not find active client connection %v", wsConn))
				continue
			}
			err = wsConn.ClientConn.WriteMessage(websocket.TextMessage, wsMsgJSON)
			if err != nil {
				svc.Log.Debug(fmt.Sprintf("Could not write to websocket %v", err))
				continue
			}
			svc.Log.Debug(fmt.Sprintf("KEEP LOOKING"))

		case commonwebsockets.Ping:
			svc.Log.Debug(fmt.Sprintf("Recieved ping %s", string(wsMsg.Uuid)))
			err := commonwebsockets.ConnectionPing(wsMsg.Uuid)
			if err != nil {
				svc.Log.Debug(fmt.Sprintf("Could not find conn for %s", string(wsMsg.Uuid)))
				continue
			}
		case commonwebsockets.Error:
			svc.Log.Debug(fmt.Sprintf("Recieved error %s", wsMsg.Data))
		case commonwebsockets.Close:
			svc.Log.Debug(fmt.Sprintf("Recieved close %s", wsMsg.Uuid))

		}

	}
	return nil
}
