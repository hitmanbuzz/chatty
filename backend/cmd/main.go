package main

import (
	"log"
	"msg_app/internal/app"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error reading .env file")
	}

	chat_app := app.Init()
	chat_app.Run()
}
