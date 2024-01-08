package httpinterface

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tvanriel/db-gui/app/domain"
)

type DumpRequest struct {
	DatabaseName        string `json:"databaseName" binding:"required"`
	CreateTableStrategy string `json:"create_table_strategy" binding:"omitempty"`
}

func DumpHandler(pool domain.ConnectionPool) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Get the connection from the pool using the token from the request
		conn := pool.Get(ctx.GetString("token"))

		// Bind the request to a DumpRequest struct
		request := &DumpRequest{}

		if err := ctx.Bind(request); err != nil {
			// Return a validation error if the request is invalid
			ctx.JSON(ValidationError(err))
			return
		}

		// Get the database from the connection
		db := conn.GetDatabase(request.DatabaseName)

		// Connect to the database
		err := db.Connect()

		if err != nil {
			// Return a server error if the connection fails
			ctx.JSON(InternalServerError(err))
			return
		}

		// Close the connection to the database
		defer db.Close()

		// Dump the database to a file
		gzFileName, err := db.Dump(request.CreateTableStrategy)

		if err != nil {
			// Return a server error if the dump fails
			ctx.JSON(InternalServerError(err))
			return
		}

		// Send the dumped file as an attachment to the client
		ctx.FileAttachment(gzFileName, getDumpFilename(request.DatabaseName))
	}
}

// getDumpFilename generates a filename for a MySQL dump file based on the
// database name and the current time.
func getDumpFilename(databaseName string) string {

	// Format the current time using the time.Stamp layout
	timestamp := time.Now().Format(time.Stamp)

	// Create and return the dump filename
	return strings.Join([]string{
		"mysqldump-",
		databaseName,
		"-",
		timestamp,
		".sql.gz",
	}, "")
}
