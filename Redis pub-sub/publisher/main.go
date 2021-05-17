package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
)

// User is a struct representing newly registered users
type User struct {
	name	 string
	location string
	gender   string
	age  	 string
	vaccine_type string
}

// MarshalBinary encodes the struct into a binary blob
// Here I cheat and use regular json :)
func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

// UnmarshalBinary decodes the struct into a User
func (u *User) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	return nil
}

// Names Some Non-Random name lists used to generate Random Users
var Names []string = []string{"Jasper", "Johan", "Edward", "Niel", "Percy", "Adam", "Grape", "Sam", "Redis", "Jennifer", "Jessica", "Angelica", "Amber", "Watch"}

// SirNames Some Non-Random name lists used to generate Random Users
var SirNames []string = []string{"Ericsson", "Redisson", "Edisson", "Tesla", "Bolmer", "Andersson", "Sword", "Fish", "Coder"}

// EmailProviders Some Non-Random email lists used to generate Random Users
var EmailProviders []string = []string{"Hotmail.com", "Gmail.com", "Awesomeness.com", "Redis.com"}

func sendMessage(msg string) {
	// Create a new Redis Client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",  // We connect to host redis, thats what the hostname of the redis service is set to in the docker-compose
		Password: "superSecret", // The password IF set in the redis Config file
		DB:       0,
	})
	// Ping the Redis server and check if any errors occured
	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		// Sleep for 3 seconds and wait for Redis to initialize
		time.Sleep(3 * time.Second)
		err := redisClient.Ping(context.Background()).Err()
		if err != nil {
			panic(err)
		}
	}
	// Generate a new background context that  we will use
	ctx := context.Background()
	// Loop and randomly generate users on a random timer
	//for {
		// Publish a generated user to the new_users channel
		//err := redisClient.Publish(ctx, "new_users", GenerateRandomUser()).Err()
		error := redisClient.Publish(ctx, "new_users", GenerateRandomUser()).Err()
		if error != nil {
			panic(error)
		}
		// Sleep random time
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(4)
		time.Sleep(time.Duration(n) * time.Second)
	//}

}

// GenerateRandomUser creates a random user, dont care too much about this.
func GenerateRandomUser(msg string) *User {
	/*
	rand.Seed(time.Now().UnixNano())
	nameMax := len(Names)
	sirNameMax := len(SirNames)
	emailProviderMax := len(EmailProviders)

	nameIndex := rand.Intn(nameMax-1) + 1
	sirNameIndex := rand.Intn(sirNameMax-1) + 1
	emailIndex := rand.Intn(emailProviderMax-1) + 1*/

	message, err := json.Marshal(msg)
	if err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	fmt.Println(string(message))

	return &User{
		name:	message.name,
		location: message.location,
		gender: message.gender,
		age: message.age,
		vaccine_type: message.vaccine_type,
		/*Username: Names[nameIndex] + " " + SirNames[sirNameIndex],
		Email:    Names[nameIndex] + SirNames[sirNameIndex] + "@" + EmailProviders[emailIndex],*/
	}
}

func http_server(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }

    switch r.Method {
		case "GET":     
			//http.ServeFile(w, r, "form.html")
			/*w.WriteHeader(http.StatusOK)
			w.Write([]byte("{\"message\": \"ok gRPC\"}"))*/
			return

		case "POST":
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}

			// Obtener el nombre enviado desde la forma
			//name := r.FormValue("name")
			// Obtener el mensaje enviado desde la forma
			msg := r.FormValue("msg")
			fmt.Println(string(msg))

			/*var body map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&body)
			failOnError(err, "Parsing JSON")
			body["origen"] = "PubSub"*/

			
			message, err := json.Marshal(msg)
			//message, err := json.Marshal(Message{name: "asdf" , location:"loc", age: "23", infectedtype: "asd", state: "as", origen: "w4" })
			//message, err := json.Marshal(body)
			// Existio un error generando el objeto JSON
			if err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}

			fmt.Println(string(message))

			sendMessage(string(message))

			fmt.Fprintf(w, "¡Mensaje Publicado!\n")
			fmt.Fprintf(w, "Message = %s\n", message)
			fmt.Fprintln(w, string(message))
		
		default:
			fmt.Fprintf(w, "Metodo %s no soportado \n", r.Method)
			return
    }
}


func main(){
	fmt.Println("Server Redis PubSub iniciado")

	http.HandleFunc("/", http_server)

	http_port := ":4000"
	
    if err := http.ListenAndServe(http_port, nil); err != nil {
        log.Fatal(err)
    }
}