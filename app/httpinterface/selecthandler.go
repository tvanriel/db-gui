package httpinterface

import (
	"github.com/gin-gonic/gin"
	"github.com/tvanriel/db-gui/app/domain"
	"github.com/tvanriel/db-gui/app/formatting"
)

func SelectHandler(pool domain.ConnectionPool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &SelectRequestConfig{}
		if err := ctx.BindJSON(request); err != nil {
			ctx.JSON(ValidationError(err))
			return
		}

		conn := pool.Get(ctx.GetString("token"))

		db := conn.GetDatabase(request.DatabaseName())

		err := db.Connect()
		if err != nil {
			ctx.JSON(InternalServerError(err))
			return
		}

		defer db.Close()

		result, sql, err := db.Select(request)

		if err != nil {
			ctx.JSON(InternalServerError(err))
			return
		}

		ctx.JSON(200, SelectResponse{
			Sql:     sql,
			Results: formatting.ToJSON(result),
		})
	}
}
