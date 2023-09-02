package main

import (
	"net/http"

	"os"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func envVariable(key string) string {

	return os.Getenv(key)
}

type jwtCustomClaims struct {
	username string `json:"username"`
	password string `json:"password"`
	Admin    bool   `json:"admin"`
	jwt.RegisteredClaims
}

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Throws unauthorized error
	if username != "jon" || password != "shhh!" {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &jwtCustomClaims{}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(envVariable("secret-key")))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func restricted(c echo.Context) error {
	if false {
		return echo.ErrBadRequest
	}

	return c.HTML(200, "<h1>wowow</h1>")
}

func register(c echo.Context) error {
	if false {
		return echo.ErrBadRequest
	}

	return c.HTML(200, "<h1>wowow</h1>")
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/login", login)

	e.GET("/register", register)

	// Restricted group
	r := e.Group("/restricted")

	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte(envVariable("secret-key")),
	}

	r.Use(echojwt.WithConfig(config))
	r.GET("/auth", restricted)

	e.Logger.Fatal(e.Start(":8081"))
}
