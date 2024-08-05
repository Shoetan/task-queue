package utils

import (
	"encoding/json"
	"log"
	"os"
	"context"
  "time"


	"github.com/go-gomail/gomail"
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



func dequeueEmailTask(client *redis.Client) (*EmailTask, error) {
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

func SendEmail(task EmailTask) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "queue@gmail.com")
	m.SetHeader("To", task.Recipent)
	m.SetHeader("Subject", task.Subject)
	m.SetBody("text/plain", task.Body)

	d := gomail.Dialer{Host: "localhost", Port: 587}

	return d.DialAndSend(m)

}

func Worker(client *redis.Client) {
	for {
		task, err := dequeueEmailTask(client)

		if err == redis.Nil {
			time.Sleep(1 * time.Second)
			continue
		} else if err != nil {
			log.Fatalf("Error dequeing task: %v\n", err)
			continue
		}

		err = SendEmail(*task)

		if err != nil {
			log.Fatalf("Error sending email:%v\n", err)
		}
	}
}