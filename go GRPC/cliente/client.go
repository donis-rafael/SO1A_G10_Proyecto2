// Paquete principal, acá iniciará la ejecución
package main

// Importar dependencias, notar que estamos en un módulo llamado tuiterclient
import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"

	"fmt"
	"io/ioutil"
	"log"
	"os"

	"net/http"

	"clientgrpc/greet.pb/greetpb"

	"google.golang.org/grpc"
)

type server struct{}

// Funcion que realiza una llamada unaria
func sendMessage(name string, location string, age string, infectedtype string, state string) {
	//server_host := os.Getenv("SERVER_HOST")
	server_host := "localhost:50051"

	fmt.Println(">> CLIENT: Iniciando cliente")
	fmt.Println(">> CLIENT: Iniciando conexion con el servidor gRPC ", server_host)

	// Crear una conexion con el servidor (que esta corriendo en el puerto 50051)
	// grpc.WithInsecure nos permite realizar una conexion sin tener que utilizar SSL
	cc, err := grpc.Dial(server_host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf(">> CLIENT: Error inicializando la conexion con el server %v", err)
	}

	// Defer realiza una axion al final de la ejecucion (en este caso, desconectar la conexion)
	defer cc.Close()

	// Iniciar un servicio NewGreetServiceClient obtenido del codigo que genero el protofile
	// Esto crea un cliente con el cual podemos escuchar
	// Le enviamos como parametro el Dial de gRPC
	c := greetpb.NewGreetServiceClient(cc)

	fmt.Println(">> CLIENT: Iniciando llamada a Unary RPC")

	// Crear una llamada de GreetRequest
	// Este codigo lo obtenemos desde el archivo que generamos con protofile
	req := &greetpb.GreetRequest{
		// Enviar un Greeting
		// Esta estructura la obtenemos desde el archivo que generamos con protofile
		Greeting: &greetpb.Greeting{
			Name:         name,
			Location:     location,
			Age:          age,
			Infectedtype: infectedtype,
			State:        state,
		},
	}

	fmt.Println(">> CLIENT: Enviando datos al server")
	// Iniciar un greet, en background con la peticion que estamos realizando
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf(">> CLIENT: Error realizando la peticion %v", err)
	}

	fmt.Println(">> CLIENT: El servidor nos respondio con el siguiente mensaje: ", res.Result)
}

// Creamos un server sencillo que unicamente acepte peticiones GET y POST a '/'
func http_server(w http.ResponseWriter, r *http.Request) {
	type Person struct {
		Name         string `json:"name"`
		Location     string `json:"location"`
		Gerder       string `json:"gender"`
		Age          int    `json:"age"`
		Vaccine_type string `json:"vaccine_type"`
	}

	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {

	// Devolver una página sencilla con una forma html para enviar un mensaje
	case "GET":
		fmt.Println(">> CLIENT: Devolviendo form.html")
		// Leer y devolver el archivo form.html contenido en la carpeta del proyecto
		http.ServeFile(w, r, "form.html")

	// Publicar un mensaje a Google PubSub
	case "POST":
		fmt.Println(">> CLIENT: Iniciando envio de mensajes")

		//Si existe un error con la forma enviada entonces no seguir
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		body, errs := ioutil.ReadAll(r.Body)
		if errs != nil {
			log.Printf("Error reading body: %v", errs)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		myBody := ioutil.NopCloser(bytes.NewBuffer(body))

		var p3 Person
		dec := json.NewDecoder(myBody)
		dec.DisallowUnknownFields()
		erre := dec.Decode(&p3)

		if erre != nil {

		}
		// Obtener el nombre enviado desde la forma
		name := p3.Name //r.FormValue("name")
		// Obtener el mensaje enviado desde la forma
		location := p3.Location //r.FormValue("location")
		age := strconv.Itoa(p3.Age)
		infectedtype := p3.Vaccine_type
		state := p3.Gerder

		// Publicar el mensaje, convertimos el objeto JSON a String
		sendMessage(name, location, age, infectedtype, state)

		// Enviamos informacion de vuelta, indicando que fue generada la peticion
		fmt.Fprintf(w, "¡Mensaje Publicado!\n")
		fmt.Fprintf(w, "Name = %s\n", name)
		fmt.Fprintf(w, "Location = %s\n", location)
		fmt.Fprintf(w, "Age = %s\n", age)
		//fmt.Fprintf(w, "Type = %s\n", infectedtype)
		//fmt.Fprintf(w, "State = %s\n", state)

	// Cualquier otro metodo no sera soportado
	default:
		fmt.Fprintf(w, "Metodo %s no soportado \n", r.Method)
		return
	}
}

// Funcion principal
func main() {
	instance_name := os.Getenv("NAME")
	client_host := os.Getenv("CLIENT_HOST")

	fmt.Println(">> -------- CLIENTE ", instance_name, " --------")

	fmt.Println(">> CLIENT: Iniciando servidor http en ", client_host)

	// Asignar la funcion que controlara las llamadas http
	http.HandleFunc("/", http_server)

	// Levantar el server, si existe un error levantandolo hay que apagarlo
	if err := http.ListenAndServe(client_host, nil); err != nil {
		log.Fatal(err)
	}
}
