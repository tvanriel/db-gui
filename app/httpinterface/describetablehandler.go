package httpinterface

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/tvanriel/db-gui/app/domain"
)

type DescribeTableResponseColumn struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Null    bool   `json:"null"`
	Default string `json:"default"`
	Comment string `json:"comment"`
}

type DescribeTableResponse struct {
	Name    string                         `json:"name"`
	Columns []*DescribeTableResponseColumn `json:"columns"`
}

func describeTableHandler(pool domain.ConnectionPool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetString("token")
		databaseName := ctx.Query("database")
		tableName := ctx.Query("table")

		conn := pool.Get(token)

		db := conn.GetDatabase(databaseName)

		if db == nil {
			ctx.JSON(ValidationError(errors.New("invalid database")))
			return
		}
		err := db.Connect()

		if err != nil {
			ctx.JSON(InternalServerError(err))
			return
		}

		defer db.Close()

		tableDescriptor, err := db.DescribeTable(tableName)

		if err != nil {
			ctx.JSON(InternalServerError(err))
			return
		}
		colNames := tableDescriptor.Columns()

		columns := []*DescribeTableResponseColumn{}
		for _, colName := range colNames {
			col := tableDescriptor.GetColumn(colName)
			columns = append(columns, &DescribeTableResponseColumn{
				Name:    col.Field(),
				Type:    col.Type(),
				Null:    col.Null(),
				Comment: col.Comment(),
				Default: col.Default(),
			})
		}

		resp := &DescribeTableResponse{
			Name:    tableName,
			Columns: columns,
		}

		ctx.JSON(200, resp)

	}

}
