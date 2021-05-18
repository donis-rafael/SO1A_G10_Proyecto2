import json
from random import random, randrange
from sys import getsizeof
from locust import HttpUser, task, between

# Esta variable controlara si queremos que salgan todas las salidas, o unicamente las mas importantes
debug = True

# Esta funcion utilizaremos para las salidas que no queremos que salgan siempre
# excepto cuando estamos debuggeando
def printDebug(msg):
    if debug:
        print(msg)

class Reader():

    def __init__(self):
        self.array = []
        
    # NOTA: ESTO QUITA EL VALOR DEL ARRAY.
    def pickRandom(self):
        length = len(self.array)
        
        if (length > 0):
            random_index = randrange(0, length - 1) if length > 1 else 0

            return self.array.pop(random_index)

        else:
            print (">> Reader: No hay más valores para leer en el archivo.")
            return None
    
    def load(self):
        print (">> Reader: Iniciando con la carga de datos")

        try:
            with open("traffic.json", 'r') as data_file:
                self.array = json.loads(data_file.read())
            
            print (f'>> Reader: Datos cargados correctamente, {len(self.array)} datos -> {getsizeof(self.array)} bytes.')

        except Exception as e:
            print (f'>> Reader: No se cargaron los datos {e}')

# Deriva de HTTP-User, simulando un usuario utilizando nuestra APP.
# En esta clase definimos todo lo que necesitamos hacer con locust.
class MessageTraffic(HttpUser):
    # Tiempo de espera entre peticiones
    #  entre cada llamada HTTP
    wait_time = between(0.1, 0.9)

    # Este metodo se ejecutara cada vez que empecemos una prueba
    # Este metodo se ejecutara POR USUARIO
    def on_start(self):
        print (">> MessageTraffic: Iniciando el envio de tráfico")
        self.reader = Reader()
        self.reader.load()

    # Este es una de las tareas que se ejecutara cada vez que pase el tiempo wait_time
    @task
    def PostMessage(self):
        random_data = self.reader.pickRandom()
        
        if (random_data is not None):

            data_to_send = json.dumps(random_data)
            printDebug (data_to_send)

            myheaders = {'Content-Type': 'application/json', 'Accept': 'application/json'}
            self.client.post("/", data= json.dumps(random_data), headers = myheaders)
            #self.client.post("/", json=random_data)

        else:
            print(">> MessageTraffic: Envio de tráfico finalizado, no hay más datos que enviar.")
            self.stop(True)

    # Este es una de las tareas que se ejecutara cada vez que pase el tiempo wait_time
    #@task
    #def GetMessages(self):
        # Realizar una peticion para recibir los datos que hemos guardado
        #self.client.get("/")