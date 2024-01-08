package httpinterface

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tvanriel/db-gui/app/domain"
)

type DbGuiHandler struct {
	engine *gin.Engine
}

func (dbh *DbGuiHandler) Handler() http.Handler {
	return dbh.engine.Handler()
}

func NewDbGuiHandler(pool domain.ConnectionPool, frontend http.FileSystem) *DbGuiHandler {
	engine := gin.Default()

	routes(pool, engine, frontend)

	return &DbGuiHandler{
		engine: engine,
	}
}

func routes(pool domain.ConnectionPool, engine *gin.Engine, frontend http.FileSystem) {

	engine.GET("/", func(ctx *gin.Context) { ctx.Redirect(301, "/app") })
	engine.StaticFS("app", frontend)

	apiV1 := engine.Group("/api/v1")
	{
		apiV1.POST(
			"/sessions",
			createSessionHandler(pool),
		)
		apiV1.POST(
			"/sessions/script",
			createSessionScriptHandler(pool),
		)

		wsAuth := apiV1.Group("", validWsTokenMiddleware(pool))
		{
			wsAuth.GET(
				"/sql",
				SqlHandler(pool),
			)
		}

		auth := apiV1.Group("", validTokenGetQueryMiddleware(pool))
		{
			auth.GET(
				"/schemas",
				listDatabasesHandler(pool),
			)

			auth.GET(
				"/tables",
				listTablesHandler(pool),
			)

			auth.GET(
				"/tables/describe",
				describeTableHandler(pool),
			)

			auth.POST(
				"/select",
				SelectHandler(pool),
			)

			auth.POST(
				"/import",
				SqlImportHandler(pool),
			)

			auth.POST("/export",
				DumpHandler(pool),
			)
		}
	}
}
