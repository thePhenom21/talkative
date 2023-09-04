package main

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/ServiceWeaver/weaver"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type jwtCustomClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}

type Auth interface {
	Login(context.Context, string, string) (string, error)
	Register(context.Context, string, string) (string, error)
}

type auth struct {
	weaver.Implements[Auth]
}

func EnvVariable(key string) string {
	godotenv.Load()
	return os.Getenv(key)
}

func (*auth) Login(c context.Context, username string, password string) (string, error) {

	path := EnvVariable("DB")

	db, err := sql.Open("postgres", path)
	if err != nil {
		print("Err" + err.Error())
	}

	q := "SELECT tok FROM users WHERE username=$1 AND password=$2 "

	row, err := db.Query(q, username, password)

	if err != nil {
		println("error " + err.Error())
	}

	if row == nil {
		return "User not found", err
	}

	if !row.Next() {
		return "User not found", err
	}

	var str_key string

	row.Scan(&str_key)

	db.Close()

	return str_key, err

}

func (*auth) Register(c context.Context, username string, password string) (string, error) {

	path := EnvVariable("DB")

	db, err := sql.Open("postgres", path)
	if err != nil {
		print("Err" + err.Error())
	}

	var checkUser string

	row, err := db.Query("SELECT username FROM users")

	if err != nil {
		println("error " + err.Error())
	}

	if row == nil {
		return "DB ERROR", err
	}

	for row.Next() {
		row.Scan(&checkUser)
		if checkUser == username {
			return "User already registered", err
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

	t, err := token.SignedString([]byte(EnvVariable("SECRET")))

	if err != nil {
		print(t)
		return "Error", err
	}

	q := "INSERT INTO users (username,password,tok) VALUES($1,$2,$3)"
	a, err := db.Exec(q, username, password, t)

	if err != nil {
		print("err" + err.Error())
		print(a.LastInsertId())
	}

	db.Close()

	return ("User " + username + " has been registerd"), err
}
