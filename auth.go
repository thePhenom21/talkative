package main

import (
	"database/sql"
	"github.com/ServiceWeaver/weaver"
	echojwt "github.com/labstack/echo-jwt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

type Auth interface {
	envVariable(string) string
	validateJWT(echo.Context) error
	login(echo.Context, *sql.DB) error
	register(echo.Context, *sql.DB) error
	init()
}

type auth struct {
	weaver.Implements[Auth]
	jwtCustomClaims
}

type jwtCustomClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}

func (a *auth) envVariable(key string) string {
	godotenv.Load()
	return os.Getenv(key)
}

func (a *auth) validateJWT(c echo.Context) error {
	return c.String(200, "hey")
}

func (a *auth) login(c echo.Context, db *sql.DB) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	q := "SELECT tok FROM users WHERE username=$1 AND password=$2 "

	row, err := db.Query(q, username, password)

	if err != nil {
		println("error " + err.Error())
	}

	if row == nil {
		return c.String(400, "User not found")
	}

	if !row.Next() {
		return c.String(400, "User not found")
	}

	var str_key string

	row.Scan(&str_key)

	return c.JSON(
		200,
		map[string]any{"key": str_key},
	)
}

func (b *auth) register(c echo.Context, db *sql.DB) error {

	username := c.QueryParam("username")
	password := c.QueryParam("password")

	var checkUser string

	row, err := db.Query("SELECT username FROM users")

	if err != nil {
		println("error " + err.Error())
	}

	if row == nil {
		return c.String(400, "db error")
	}

	for row.Next() {
		row.Scan(&checkUser)
		if checkUser == username {
			return c.String(400, "User already registered")
		}
	}

	claims := &jwtCustomClaims{
		username,
		password,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(b.envVariable("SECRET")))

	if err != nil {
		print(t)
		return err
	}

	q := "INSERT INTO users (username,password,tok) VALUES($1,$2,$3)"
	a, err := db.Exec(q, username, password, t)

	if err != nil {
		print("err" + err.Error())
		print(a.LastInsertId())
	}

	return c.String(200, "User "+username+" has been registered")
}

func (a *auth) init() {
	e := echo.New()

	path := a.envVariable("DB")

	db, err := sql.Open("postgres", path)
	if err != nil {
		print("Err" + err.Error())
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	r := e.Group("/restricted")

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte(a.envVariable("SECRET")),
	}
	r.Use(echojwt.WithConfig(config))

	e.POST("/login", func(c echo.Context) error {
		return a.login(c, db)
	})

	e.POST("/register", func(c echo.Context) error {
		return a.register(c, db)
	})

	r.GET("", a.validateJWT)

	defer db.Close()

	e.Logger.Fatal(e.Start(":8081"))
}
