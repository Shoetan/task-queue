package main

import (
	//"github.com/Shoetan/pkg/db"
	"github.com/Shoetan/pkg/server"
)


func main()  {
	// db.Database()
	// start the server 

	server := server.NewAPISERVER(":4005")

	server.Run()
}