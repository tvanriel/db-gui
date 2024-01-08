package scriptauthentication

import (
	"context"
	"time"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
)

type AuthScriptResponse struct {
	Username string
	Password string
	Expires  int
	Host     string
}

func ExecuteAuthScript(scriptSrc string) (*AuthScriptResponse, error) {
	script := tengo.NewScript([]byte(scriptSrc))
	script.SetImports(stdlib.GetModuleMap(stdlib.AllModuleNames()...))

	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	compiled, err := script.RunContext(ctx)
	if err != nil {
		return nil, err
	}

	hostValue := compiled.Get("host").String()
	userValue := compiled.Get("username").String()
	passwordValue := compiled.Get("password")
	expiresValue := compiled.Get("expires")
	scriptError := compiled.Get("err")

	if scriptError.Error() != nil {
		return nil, scriptError.Error()
	}

	return &AuthScriptResponse{
		Username: userValue,
		Password: passwordValue.String(),
		Host:     hostValue,
		Expires:  expiresValue.Int(),
	}, nil
}
