package main

import (
	// "github.com/aws/aws-lambda-go/lambda"

	"valorant/services"
  "valorant/db"
)

func update_database() {
  db.AutoMigrate()
  users := services.GetUsers()
  for _, user := range(users) {
    db.InsertMMRHistory(user)
  }
}

func main() {
  // lambda.Start(update_database)
  update_database()
}
