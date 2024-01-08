package httpinterface

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/tvanriel/db-gui/app/domain"
)

// validTokenGetQueryMiddleware is a middleware function that checks the validity of a token passed in the query string.
// It takes a ConnectionPool as an argument and returns a gin.HandlerFunc.
func validTokenGetQueryMiddleware(pool domain.ConnectionPool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the token from the query string
		token := ctx.Query("token")

		// Get the connection associated with the token from the connection pool
		conn := pool.Get(token)

		if conn == nil {
			// If the connection is not found, return a 401 Unauthorized response
			ctx.JSON(Unauthorized("Invalid token"))
			ctx.Abort()
			return
		}
		// Set the token in the request context
		ctx.Set("token", token)
	}
}

type WsError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// validWsTokenMiddleware is a middleware function that upgrades an HTTP request to a WebSocket connection and checks the validity of a token passed in the query string.
// It takes a ConnectionPool as an argument and returns a gin.HandlerFunc.
func validWsTokenMiddleware(pool domain.ConnectionPool) gin.HandlerFunc {
	// Create a websocket.Upgrader with some default values
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// Allow all origins
		CheckOrigin: func(*http.Request) bool { return true },
	}

	return func(ctx *gin.Context) {
		// Upgrade the HTTP request to a WebSocket connection
		ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			return
		}

		// Get the token from the query string
		token := ctx.Query("token")

		// Get the connection associated with the token from the connection pool
		conn := pool.Get(token)

		if conn == nil {
			// If the connection is not found, send an error message over the WebSocket connection and return
			ws.WriteJSON(WsError{
				Type:    "sql/handler/error",
				Message: "Invalid token",
			})
			ws.Close()
			ctx.Abort()
			return
		}

		// Set the token and WebSocket connection in the request context
		ctx.Set("token", token)
		ctx.Set("ws", ws)
	}
}
