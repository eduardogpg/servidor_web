package main

import(
	"github.com/gorilla/mux"
	"net/http"
	"log"
)

func HomeHandler(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Hola Mundo\n"))
}

func main() {
	mux := mux.NewRouter()
	mux.HandleFunc("/", HomeHandler)
	
	http.Handle("/", mux)
	log.Println("El servidor se encuentra a la escucha en el puerto 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
	//Es equivalente al println solo que este terminame el programa
}