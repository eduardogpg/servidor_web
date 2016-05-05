package main

import(
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"encoding/json"
)

type Response struct {
	Message string `json:"user_name"`
}

func HolaMundo(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Hola Mundo\n"))
}

func HolaJson(w http.ResponseWriter, r *http.Request){
	response := Response{"Estes es un mensaje en JSON"}
	json.NewEncoder(w).Encode(response)
}

func HomeHandler(w http.ResponseWriter, r *http.Request){
	
}

func main() {
	mux := mux.NewRouter()
	mux.HandleFunc("/", HolaMundo)
	mux.HandleFunc("/json", HolaJson)
	
	http.Handle("/", mux)
	log.Println("El servidor se encuentra a la escucha en el puerto 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
	//Es equivalente al println solo que este terminame el programa
}