package httpinterface

import (
	"github.com/gin-gonic/gin"
	"github.com/tvanriel/db-gui/app/domain"
)

func listDatabasesHandler(pool domain.ConnectionPool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetString("token")
		conn := pool.Get(token)

		err := conn.Connect()

		if err != nil {
			ctx.JSON(InternalServerError(err))
		}

		databases, err := conn.GetDatabases()

		if err != nil {
			ctx.JSON(InternalServerError(err))
			return
		}

		conn.Close()

		type listSchemaResponseItem struct {
			Item string `json:"item"`
		}

		response := []listSchemaResponseItem{}

		for _, database := range databases {
			response = append(response, listSchemaResponseItem{Item: database})
		}

		ctx.JSON(200, response)

	}
}
