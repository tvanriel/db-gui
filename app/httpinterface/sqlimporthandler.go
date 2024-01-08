package httpinterface

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	archiver "github.com/mholt/archiver/v3"
	"github.com/tvanriel/db-gui/app/communicator"
	"github.com/tvanriel/db-gui/app/domain"
	"github.com/tvanriel/db-gui/app/mysqlconnection/queries"
)

type ImportError struct {
	Error string  `json:"error"`
	Sql   *string `json:"sql"`
}

type SqlImportResponse struct {
	Errors   []ImportError `json:"errors"`
	Executed int           `json:"executed"`
}

func SqlImportHandler(pool domain.ConnectionPool) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Parse the request
		req, err := c.MultipartForm()
		if err != nil {
			c.AbortWithStatusJSON(ValidationError(err))
			return
		}

		// Get the file from the request
		file, ok := req.File["file"]
		if !ok {
			c.AbortWithStatusJSON(ValidationError(errors.New("missing file")))
			return
		}

		databaseName := ""
		databaseNameField, ok := req.Value["databaseName"]

		if ok && len(databaseNameField) > 0 {
			databaseName = databaseNameField[0]
		}

		// Read the file into a byte slice
		data, err := file[0].Open()
		if err != nil {
			c.AbortWithStatusJSON(InternalServerError(err))
			return
		}

		defer data.Close()

		buf := bytes.NewBuffer(make([]byte, 0))

		// Check if the file is a Gzip archive
		if err = archiver.DefaultGz.CheckExt(file[0].Filename); err == nil {
			// Extract the SQL file from the Gzip archive
			err := archiver.DefaultGz.Decompress(data, buf)
			if err != nil {
				c.AbortWithStatusJSON(InternalServerError(err))
				return
			}
		} else {
			io.Copy(buf, data)
		}

		statements, err := queries.ParseScript(buf.String(), func(i int) {})

		if err != nil {
			c.JSON(InternalServerError(err))
			return
		}
		errors := []ImportError{}

		conn := pool.Get(c.GetString("token"))

		q, err := getQueriable(conn, databaseName)
		if err != nil {
			c.JSON(InternalServerError(err))
			return
		}

		communicator := communicator.NewErrorCommunicator(
			func(err error, sql string) {
				errors = append(errors, ImportError{
					Error: err.Error(),
					Sql:   &sql,
				})
			},
			func(err error) {
				errors = append(errors, ImportError{
					Error: err.Error(),
				})
			},
		)
		executeSqlScript(statements, communicator, q)

		c.JSON(http.StatusOK, SqlImportResponse{
			Errors:   errors,
			Executed: len(statements),
		})

	}
}
