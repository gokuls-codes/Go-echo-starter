package server

import (
	"database/sql"

	"github.com/gokuls-codes/go-echo-starter/internal/services/users"
	"github.com/gokuls-codes/go-echo-starter/templates/pages"
	"github.com/gokuls-codes/go-echo-starter/types"
	"github.com/gokuls-codes/go-echo-starter/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	customMiddleware "github.com/gokuls-codes/go-echo-starter/internal/middleware"
)

type Server struct {
	addr string
	db *sql.DB
}

func NewServer(addr string, db *sql.DB) *Server {
	return &Server{
		addr: addr,
		db: db,
	}
}

func (s *Server) Start() error {
	app := echo.New()

	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	app.Static("/static", "assets")

	app.GET("/robots.txt", func(c echo.Context) error {
		return c.File("robots.txt")
	})

	app.Use(customMiddleware.Theme)

	userGroup := app.Group("/auth")
	userStore := users.NewStore(s.db)
	userHandler := users.NewHandler(userStore)
	userHandler.RegisterRoutes(userGroup)

	homeGroup := app.Group("")
	homeGroup.Use(customMiddleware.Auth(userStore))

	homeGroup.GET("/", func(c echo.Context) error {
		return utils.Render(c, pages.HomePage(c.Get("theme")=="dark", c.Get("user").(*types.User).Name))
	})

	return app.Start(s.addr)
}