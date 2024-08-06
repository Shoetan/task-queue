package main

import (

	"github.com/Shoetan/pkg/server"
	"github.com/Shoetan/pkg/utils"
)


func main()  {
	// start the server 

	redisClient := utils.RedisClient()

	go utils.Worker(redisClient)
	
	server := server.NewAPISERVER(":4005")

	server.Run()


	




}