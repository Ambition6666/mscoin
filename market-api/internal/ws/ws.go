package ws

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"strings"
)

const ROOM = "market"

type WebSocketHandler func(s socketio.Conn) error

type WebSocketServer struct {
	wsServer   *socketio.Server
	pathRouter httpx.Router
	path       string
}

func (w *WebSocketServer) Start() {
	logx.Info("============socketIO启动================")
	w.wsServer.Serve()
}

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

func NewWebSocketServer(r httpx.Router, path string) *WebSocketServer {
	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})
	w := &WebSocketServer{
		pathRouter: r,
		wsServer:   server,
		path:       path,
	}
	w.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		s.Join(ROOM)
		return nil
	})
	return w
}
func (w *WebSocketServer) Stop() {
	logx.Info("============socketIO关闭================")
	w.wsServer.Close()
}

func (w *WebSocketServer) OnConnect(path string, handler WebSocketHandler) {
	w.wsServer.OnConnect(path, handler)
}
func (w *WebSocketServer) BroadcastToNamespace(path string, event string, data any) {
	go func() {
		w.wsServer.BroadcastToRoom(path, ROOM, event, data)
	}()
}

func (ws *WebSocketServer) ServerHandler(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logx.Info("=============", r.URL.Path)
		if strings.HasPrefix(r.URL.Path, ws.path) {
			ws.wsServer.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
