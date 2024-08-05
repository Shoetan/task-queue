package utils

import (
	"encoding/json"
	"log"
	"os"
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)


type EmailTask struct {
	Recipent string
	Subject string
	Body string
}

var ctx = context.Background()


func GetEnvVariable(key string) string{
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Unable to load environmental variable ")
	}
	return os.Getenv(key)
}


func HashPassword(password string) (string, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPwd), err

}

// create the task queue using the redis  

func EnqueueEmailTask(client *redis.Client, task EmailTask) error {
	jsonTask, err := json.Marshal(task)

	if err != nil {
		return err
	}
	return client.RPush(ctx, "email_task", jsonTask).Err()
}



func DequeueEmailTask(client *redis.Client) (*EmailTask, error) {
	jsonTask, err := client.LPop(ctx, "email_task").Result()

	if err != nil {
		return nil, err
	}

	var task EmailTask

	err = json.Unmarshal([]byte(jsonTask), &task)

	if err != nil {
		return nil, err
	}
	return &task, nil
}