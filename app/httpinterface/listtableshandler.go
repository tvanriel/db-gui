package httpinterface

import (
	"github.com/gin-gonic/gin"
	"github.com/tvanriel/db-gui/app/domain"
)

func listTablesHandler(pool domain.ConnectionPool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetString("token")
		databaseName := ctx.Query("database")

		conn := pool.Get(token)

		if conn == nil {
			ctx.JSON(Unauthorized("Invalid token"))
			return
		}

		db := conn.GetDatabase(databaseName)
		err := db.Connect()

		if err != nil {
			ctx.JSON(InternalServerError(err))
			return
		}

		defer db.Close()

		tables, err := db.GetTableNames()

		if err != nil {
			ctx.JSON(InternalServerError(err))
			return
		}

		ctx.JSON(200, tables)
	}
}
