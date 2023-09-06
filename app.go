package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/ServiceWeaver/weaver"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	upgrader = websocket.Upgrader{}
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
	authClient weaver.Ref[Auth]
	hl         weaver.Listener
}

// serve is called by weaver.Run and contains the body of the application.
func serve(ctx context.Context, app *app) error {
	e := echo.New()

	path := EnvVariable("DB")

	db, _ := sql.Open("postgres", path)

	_, err := db.Query("SELECT * FROM users")

	if err != nil {
		db.Exec("CREATE TABLE users(username varchar,password varchar,tok varchar)")
	}

	db.Close()

	var authClient = app.authClient.Get()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	r := e.Group("/chat")

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte(EnvVariable("SECRET")),
	}

	r.Use(echojwt.WithConfig(config))

	r.GET("", func(c echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer ws.Close()

		for {
			// Write
			err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
			if err != nil {
				c.Logger().Error(err)
			}

			// Read
			_, msg, err := ws.ReadMessage()
			if err != nil {
				c.Logger().Error(err)
			}
			fmt.Printf("%s\n", msg)
		}
	})

	e.POST("/register", func(c echo.Context) error {
		returnVal, _ := authClient.Register(ctx, c.QueryParam("username"), c.QueryParam("password"))

		return c.String(200, returnVal)

	})

	e.GET("/test", func(c echo.Context) error {
		return c.String(200, "connection succesful")
	})

	e.POST("/login", func(c echo.Context) error {
		returnVal, _ := authClient.Login(ctx, c.QueryParam("username"), c.QueryParam("password"))

		return c.String(200, returnVal)
	})

	e.Listener = app.hl

	return e.Start("")

}
