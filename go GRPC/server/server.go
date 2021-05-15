// Paquete principal, acá iniciará la ejecución
package main

// Importar dependencias, notar que estamos en un módulo llamado grpctuiter
import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gomodule/redigo/redis"

	"log"

	"servergrpc/greet.pb/greetpb"

	"google.golang.org/grpc"
)

// Iniciar una estructura que posteriormente gRPC utilizará para realizar un server
type server struct{}

// Función que será llamada desde el cliente
// Debemos pasarle un contexto donde se ejecutara la funcion
// Y utilizar las clases que fueron generadas por nuestro proto file
// Retornara una respuesta como la definimos en nuestro protofile o un error
func (*server) Greet(ctxt context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf(">> SERVER: Función Greet llamada con éxito. Datos: %v\n", req)

	// Todos los datos podemos obtenerlos desde req
	// Tendra la misma estructura que definimos en el protofile
	// Para ello utilizamos en este caso el GetGreeting
	Name := req.GetGreeting().GetName()
	Location := req.GetGreeting().GetLocation()
	Age := req.GetGreeting().GetAge()
	Infectedtype := req.GetGreeting().GetInfectedtype()
	State := req.GetGreeting().GetState()

	result := Name + " - " + Location + " - " + Age + " - " + Infectedtype + " - " + State

	fmt.Printf(">> SERVER: %s\n", result)
	// Creamos un nuevo objeto GreetResponse definido en el protofile

	jsonData := map[string]string{"name": Name, "location": Location, "age": Age, "vaccine_type": Infectedtype, "gender": State, "origen": "grpc"}
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

	insertResult, err := collection.InsertOne(context.TODO(), jsonData)

	//resp, err := client.Do(request)

	//request, err = http.Post("https://localhost/subscribers", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf(">> Error 3")
		fmt.Printf("The request failed with error %s\n", err)
	} else {
		fmt.Println("Dato insertado ID: ", insertResult.InsertedID)

		c, err := redis.Dial("tcp", "35.184.169.129:6379")
		fmt.Printf("Conectado a Redis %s\n", err)
		if err != nil {
			fmt.Printf("No se pudo conectar a Redis %s\n", err)
		} else {
			b, err := json.Marshal(jsonData)
			if err != nil {
				fmt.Printf("Error:  %s\n", err)
			}

			if _, err := c.Do("RPUSH", "sopes", string(b)); err != nil {
				fmt.Printf("No se pud insertar a Redis %s\n", err)
			}

		}
	}
	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

// Funcion principal
func main() {

	// Leer el host de las variables del ambiente
	//host := os.Getenv("HOST")
	host := "localhost:50051"
	fmt.Println(">> SERVER: Iniciando en ", host)

	// Primero abrir un puerto para poder escuchar
	// Lo abrimos en este puerto arbitrario
	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf(">> SERVER: Error inicializando el servidor: %v", err)
	}

	fmt.Println(">> SERVER: Empezando server gRPC")

	// Ahora si podemos iniciar un server de gRPC
	s := grpc.NewServer()

	// Registrar el servicio utilizando el codigo que nos genero el protofile
	greetpb.RegisterGreetServiceServer(s, &server{})

	fmt.Println(">> SERVER: Escuchando servicio...")
	// Iniciar a servir el servidor, si hay un error salirse
	if err := s.Serve(lis); err != nil {
		log.Fatalf(">> SERVER: Error inicializando el listener: %v", err)
	}
}
