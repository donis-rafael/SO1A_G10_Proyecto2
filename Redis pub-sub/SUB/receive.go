package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func fnMensajeError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func SendPostAsync(url string, body []byte, rc chan *http.Response) {
	response, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	rc <- response
}

func main() {
	type Person struct {
		Name         string `json:"name"`
		Location     string `json:"location"`
		Gerder       string `json:"gender"`
		Age          int    `json:"age"`
		Vaccine_type string `json:"vaccine_type"`
	}
	// Conectando con el servidor de redis, puerto default en redis 6379
	//opt, err := redis.ParseURL("redis://localhost:6379/0")
	address, ok := os.LookupEnv("SERVER_HOST")
	if !ok {
		address = "redis://35.184.169.129:6379/0"
	}
	log.Printf("Llamando in port %s", address)

	opt, err := redis.ParseURL(address)

	fnMensajeError(err, "Failed to connect to redis")
	redis := redis.NewClient(opt)
	redissubscribe := redis.Subscribe(context.Background(), "redisg10")
	canal := redissubscribe.Channel()
	for msg := range canal {

		log.Printf("Received a message: %s", msg.Payload)

		var persona Person
		json.Unmarshal([]byte(msg.Payload), &persona)
		//post := []byte(msg.Payload)
		//response := make(chan *http.Response)
		// go SendPostAsync("http://35.209.183.162/casos", post, response)
		//out, err := json.Marshal(msg.Payload)

		//jsonData := map[string]string{"name": Name, "location": Location, "age": Age, "vaccine_type": Infectedtype, "gender": State, "origen": "grpc"}
		//client := &http.Client{}
		//client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://sopes:sopes123@localhost:27017/subscribers").SetAuth(options.Credential{AuthSource: "subscribers", Username: "sopes", Password: "sopes123"}))
		//	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://sopes:sopes123@localhost:27017/subscribers"))
		// if err != nil {
		// 	fmt.Printf(">> Error 1")
		// 	fmt.Print(err.Error())
		// }
		//result.Header.Add("Accept", "application/json")
		//result.Header.Add("Content-Type", "application/json")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://sopes:sopes123@35.184.169.129:27017/subscribers"))

		//err = client.Connect(ctx)
		if err != nil {
			fmt.Printf(">> Error 2")
			fmt.Print(err.Error())
		}

		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				panic(err)
			}
		}()

		// Ping the primary
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			fmt.Printf("Error ping.\n")
			//log.Fatal(err)
		}
		fmt.Println("Successfully connected and pinged.")

		collection := client.Database("subscribers").Collection("subscribers")
		//jsonValue, _ := json.Marshal(jsonData)
		fmt.Printf("Successfully connected and pinged.  %s", persona.Name)
		jsonData := map[string]string{"name": persona.Name, "location": persona.Location, "age": strconv.Itoa(persona.Age), "vaccine_type": persona.Vaccine_type, "gender": persona.Gerder, "origen": "redis"}

		insertResult, err := collection.InsertOne(context.TODO(), jsonData)

		//resp, err := client.Do(request)

		//request, err = http.Post("https://localhost/subscribers", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Printf(">> Error 3")
			fmt.Printf("The request failed with error %s\n", err)
		} else {
			fmt.Println("Dato insertado ID: ", insertResult.InsertedID)

			//c, err := redis.Dial("tcp", "35.184.169.129:6379")
			// c, err := redis.NewClient(&redis.Options{
			// 	Addr: "localhost:6379",
			// })
			fmt.Printf("Conectado a Redis %s\n", err)
			if err != nil {
				fmt.Printf("No se pudo conectar a Redis %s\n", err)
			} else {
				// b, err := json.Marshal(response)
				// if err != nil {
				// 	fmt.Printf("Error:  %s\n", err)
				// }

				//log.Printf("%s: ", message)
				errs := redis.Publish(context.TODO(), "sopes", string(msg.Payload)).Err()

				if errs != nil {
					fnMensajeError(errs, "Error parseando Json")
				}

			}
		}

	}
}
