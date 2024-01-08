package httpinterface

import (
	"encoding/json"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/tvanriel/db-gui/app/domain"
)

type CreateSessionRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Host     string `json:"host" binding:"required"`
}

// createSessionHandler is a handler function that handles requests to create a new session.
func createSessionHandler(pool domain.ConnectionPool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Bind the request payload to a CreateSessionRequest struct
		request := &CreateSessionRequest{}

		if err := ctx.BindJSON(request); err != nil {
			// If the request payload is invalid, return a validation error
			ctx.AbortWithStatusJSON(ValidationError(err))
			return
		}

		// Generate a new token
		token := pool.NewToken()

		script, err := makeScriptFromSessionRequest(request)

		if err != nil {
			ctx.AbortWithStatusJSON(ValidationError(err))
			return
		}

		// Add a new connection to the connection pool using the provided parameters and the generated token
		err = pool.Add(domain.NewConnectionConfig(script), token)

		if err != nil {
			// If there is an error adding the connection to the pool, return a validation error
			ctx.AbortWithStatusJSON(ValidationError(err))
			pool.Delete(token)
			return
		}

		// Return the generated token in the response
		ctx.JSON(200, struct {
			Token string `json:"token"`
		}{Token: token})
	}
}

type CreateScriptSessionRequest struct {
	Script string `json:"script"`
}

func createSessionScriptHandler(pool domain.ConnectionPool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Bind the request payload to a CreateSessionRequest struct
		request := &CreateScriptSessionRequest{}

		if err := ctx.BindJSON(&request); err != nil {
			// If the request payload is invalid, return a validation error
			ctx.JSON(ValidationError(err))
			return
		}

		// Generate a new token
		token := pool.NewToken()

		// Add a new connection to the connection pool using the provided parameters and the generated token
		err := pool.Add(domain.NewConnectionConfig(request.Script), token)

		if err != nil {
			// If there is an error adding the connection to the pool, return a validation error
			ctx.JSON(ValidationError(err))
			pool.Delete(token)
			return
		}

		// Return the generated token in the response
		ctx.JSON(200, struct {
			Token string `json:"token"`
		}{Token: token})
	}
}

var scriptTemplate = template.Must(template.New("script_from_session_request").Parse(`username:={{.Username}}
password:={{.Password}}
host:={{.Host}}
expires:=60*60*24`))

type scriptTemplateParams struct {
	Username string
	Password string
	Host     string
}

func makeScriptFromSessionRequest(request *CreateSessionRequest) (string, error) {
	usernameEscaped, err := json.Marshal(request.Username)
	if err != nil {
		return "", err
	}
	passwordEscaped, err := json.Marshal(request.Password)
	if err != nil {
		return "", err
	}
	hostnameEscaped, err := json.Marshal(request.Host)
	if err != nil {
		return "", err
	}

	strbuf := &strings.Builder{}
	err = scriptTemplate.Execute(strbuf, scriptTemplateParams{
		Username: string(usernameEscaped),
		Host:     string(hostnameEscaped),
		Password: string(passwordEscaped),
	})
	if err != nil {
		return "", err
	}
	return strbuf.String(), err
}
