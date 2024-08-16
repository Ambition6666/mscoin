package handler

import (
	"github.com/zeromicro/go-zero/rest"
	"net/http"
)

type Routers struct {
	server      *rest.Server
	middlewares []rest.Middleware
}

func NewRouters(server *rest.Server) *Routers {
	return &Routers{server: server}
}

func (r *Routers) Get(path string, handlerFunc http.HandlerFunc) {
	r.server.AddRoutes(
		rest.WithMiddlewares(
			r.middlewares,
			rest.Route{
				Method:  http.MethodGet,
				Path:    path,
				Handler: handlerFunc,
			},
		),
	)
}
func (r *Routers) Post(path string, handlerFunc http.HandlerFunc) {
	r.server.AddRoutes(
		rest.WithMiddlewares(
			r.middlewares,
			rest.Route{
				Method:  http.MethodPost,
				Path:    path,
				Handler: handlerFunc,
			},
		),
	)
}

func (r *Routers) Group() *Routers {
	return &Routers{server: r.server}
}

func (r *Routers) Use(middlewares ...rest.Middleware) {
	r.middlewares = middlewares
}
