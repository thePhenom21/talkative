package main

import (
	"context"
	"database/sql"
	"github.com/ServiceWeaver/weaver"
	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
)

func main() {
	if err := weaver.Run[app, *app](context.Background(), serve); err != nil {
		log.Fatal(err)
	}
}

// app is the main component of the application. weaver.Run creates
// it and passes it to serve.
type app struct {
	weaver.Implements[weaver.Main]
	authClient   weaver.Ref[Auth]
	chatClient   weaver.Ref[Chat]
	httpListener weaver.Listener
}

// serve is called by weaver.Run and contains the body of the application.
func serve(ctx context.Context, app *app) error {
	e := echo.New()

	var authClient = app.authClient.Get()
	var chatClient = app.chatClient.Get()

	path := authClient.envVariable("DB")

	db, err := sql.Open("postgres", path)
	if err != nil {
		print("Err" + err.Error())
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	r := e.Group("/chat")

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte(authClient.envVariable("SECRET")),
	}
	r.Use(echojwt.WithConfig(config))

	e.POST("/login", func(c echo.Context) error {
		return authClient.login(c, db)
	})

	e.POST("/register", func(c echo.Context) error {
		return authClient.register(c, db)
	})

	r.GET("", chatClient.hello)

	defer db.Close()

	e.Listener = app.httpListener

	return e.Start("localhost:8080")

}
