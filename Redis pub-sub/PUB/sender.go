package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

// Publicando nuevo mensaje
func messge(message []byte) {
	address, ok := os.LookupEnv("SERVER_HOST")
	if !ok {
		address = "redis://35.184.169.129:6379/0"
	}
	log.Printf("Llamando in port %s", address)

	opt, err := redis.ParseURL(address)

	if err != nil {
		panic(err)
	}

	redis := redis.NewClient(opt)
	//log.Printf("%s: ", message)
	errs := redis.Publish(context.TODO(), "redisg10", message).Err()

	if errs != nil {
		fnMensajeError(errs, "Error parseando Json")
	}
}

func fnNuevoElemnto(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		requestAt := time.Now()
		//Agregr encabezado
		w.Header().Set("Content-Type", "application/json")

		// Parseando json
		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		log.Printf(" parseando json")
		fnMensajeError(err, "Error parseando Json")

		//campo source para identificar la ruta de ingreso.
		body["origen"] = "PubSub"
		data, err := json.Marshal(body)

		fnMensajeError(err, "PubSub connection")
		messge(data)
		duration := time.Since(requestAt)
		fmt.Fprintf(w, "Task scheduled in %+v", duration)

	case "GET":
		fmt.Fprintf(w, "Hello")

	default:
		http.Error(w, "Metodo no soportado", 401)
		fmt.Fprintf(w, "Metodo no soportado")
		return
	}
}

func fnMensajeError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

//escuchando en el puerto 80
func fnSolicitudes() {
	http.HandleFunc("/", fnNuevoElemnto)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func main() {
	fnSolicitudes()
}
