package httpinterface

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/tvanriel/db-gui/app/communicator"
	"github.com/tvanriel/db-gui/app/domain"
	"github.com/tvanriel/db-gui/app/mysqlconnection/queries"
)

type SqlExecuteRequest struct {
	DatabaseName string `json:"databaseName"`
	SQL          string `json:"sql" binding:"required"`
}

func SqlHandler(pool domain.ConnectionPool) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		// Get the token and database name from the request.
		token := ctx.GetString("token")

		var ws *websocket.Conn
		var ok bool
		if ws, ok = (ctx.MustGet("ws")).(*websocket.Conn); !ok {
			ctx.JSON(InternalServerError(errors.New("type-error: ws is not a websocket")))
			return
		}

		request := &SqlExecuteRequest{}
		communicator := communicator.NewWsCommunicator(ws)
		err := ws.ReadJSON(request)

		if request.SQL == "" {
			communicator.GenericError(errors.New("SQL string may not be empty"))
			communicator.Close()
			return
		}

		if err != nil {
			communicator.GenericError(err)
			communicator.Close()
			return
		}

		// Get the connection from the pool and attempt to connect
		conn := pool.Get(token)

		q, err := getQueriable(conn, request.DatabaseName)
		// Failed to make a connection.  Return it as an internal server error.
		if err != nil {
			communicator.GenericError(err)
			communicator.Close()
			return
		}

		defer q.Close()

		statements, err := queries.ParseScript(request.SQL, communicator.ParsingProgress)

		// Failed to parse script.
		if err != nil {
			communicator.GenericError(err)
			communicator.Close()
			return
		}

		executeSqlScript(statements, communicator, q)
		ws.Close()

	}
}
