package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
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

	if err == redis.Nil{
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("could not dequeue task: %v", err)
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
	m.SetHeader("From", "emmysoetan@gmail.com")
	m.SetHeader("To", task.Recipent)
	m.SetHeader("Subject", task.Subject)
	m.SetBody("text/plain", task.Body)

	d := gomail.Dialer{Host: "localhost", Port: 1025}

	return d.DialAndSend(m)

}

func Worker(client *redis.Client) {
	fmt.Println("Starting worker 🛠 ...")
	for {
		task, err := dequeueEmailTask(client)
		if err == nil && task == nil {
			time.Sleep(1 * time.Second)
			continue
		} else if err != nil {
			log.Printf("Error dequeing task: %v\n", err)
			continue
		}


		fmt.Println("Task dequeued successfully:",task)
		

		tasks, err := ListAllTasks(client)

		if err != nil {
			fmt.Printf("Error retreiving task: %v", err)
		} else {
			fmt.Println("Task in list:")
			for _, task := range tasks{
				fmt.Println(task)
			}
		}

		err = SendEmail(*task)

		if err != nil {
			log.Printf("Error sending email:%v\n", err)
		}
	}
}

func RedisClient() *redis.Client{
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return redisClient
}


func ListAllTasks(client *redis.Client) ([]string, error){
	task, err := client.LRange(ctx, "email_task", 0, -1).Result()

	if err != nil {
		return nil, fmt.Errorf("could not retrieve task: %v", err)
	}

	return task, nil

}