package main

import(
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"encoding/json"
	"sync"
)

var Users = make(map[string]User)
var UserssRWMutex sync.RWMutex 
// RWMutex is a reader/writer mutual exclusion lock
// Ver Documentación

// Maps are not safe for concurrent use: 
// it's not defined what happens when you read and write to them simultaneously.
// If you need to read from and write to a map from concurrently executing goroutines,
// the accesses must be mediated by some kind of synchronization mechanism.
// One common way to protect maps is with sync.RWMutex.

type User struct{
	Websocket *websocket.Conn
	User_Name string
}

type Response struct {
	Message string 	`json:"message"`
	Valid 	bool		`json:"valid"`
}

func HolaMundo(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Hola Mundo\n"))
}

func HolaJson(w http.ResponseWriter, r *http.Request){
	response := CreateResponse("Estes es un mensaje en JSON", true)
	json.NewEncoder(w).Encode(response)
}

func CreateResponse( message string, valid bool ) Response{
	return Response{ message, valid}
}


func HomeStaticPage(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w,r,"Front/index.html")
}

func HomeHandler(w http.ResponseWriter, r *http.Request){
}

func Validate(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	user_name := r.FormValue("user_name")
	response := Response{}
	response.Message = user_name

	if !UserExist(user_name){
		response.Valid = true;
	}else{
		response.Valid = false;
	}
	json.NewEncoder(w).Encode(response)
}

func UserExist(user_name string) bool{
	UserssRWMutex.Lock()
	defer UserssRWMutex.Unlock()

	if _, ok := Users[user_name]; ok{
		return true
	}
	return false
}

func CreateUser(user_name string, ws *websocket.Conn) User{
	return User{User_Name: user_name, Websocket: ws}
}

func AddUser(user User){
	UserssRWMutex.Lock()
	defer UserssRWMutex.Unlock()
	Users[user.User_Name] = user
}

func RemoveUser(user_name string){
	UserssRWMutex.Lock()
	defer UserssRWMutex.Unlock()
	delete(Users, user_name)
}

func SenMessageUsers(messageType int, message []byte){
	UserssRWMutex.Lock()
	defer UserssRWMutex.Unlock()

	for _, user := range Users{
		if err := user.Websocket.WriteMessage(messageType, message); err != nil{
			return
		}
	}
}

func LenMap()int{
	UserssRWMutex.Lock()
	defer UserssRWMutex.Unlock()
	return len(Users)
}

func GetArrayByte(value string)[]byte{
		return []byte(value)
}

func GetStringByte(array [] byte) string{
	return string(array[:])
}

func GetMessageFormat(user_name string, message []byte)[]byte{
	return GetArrayByte(user_name + ":" + GetStringByte(message))
}

func WebSocket(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	user_name := vars["user_name"]

	ws, err := websocket.Upgrade(w,r,nil,1024, 1024)
	//(w http.ResponseWriter, r *http.Request, responseHeader http.Header, readBufSize, writeBufSize int) (*Conn, error)
	if err != nil{
		log.Println(err)
		return
	}
	current_user := CreateUser(user_name, ws)
	AddUser(current_user)
	for{
		type_message, message, err := ws.ReadMessage()
		//messageType int, p []byte, err error
		if err != nil{
			RemoveUser(user_name)
			return
		}
		final_message := GetMessageFormat(user_name,message)
		SenMessageUsers(type_message, final_message)
	}
}


func main() {
	mux := mux.NewRouter()
	//https://golang.org/pkg/net/http/
	cssHandler := http.FileServer(http.Dir("./Front/CSS/"))
	js_Handler := http.FileServer(http.Dir("./Front/JS/"))

	mux.HandleFunc("/", HomeStaticPage).Methods("GET")
	mux.HandleFunc("/hola", HolaMundo).Methods("GET")
	mux.HandleFunc("/json", HolaJson).Methods("GET")
	mux.HandleFunc("/validate", Validate).Methods("POST")
	mux.HandleFunc("/chat/{user_name}", WebSocket).Methods("GET")

	
	http.Handle("/", mux)
	http.Handle("/CSS/", http.StripPrefix("/CSS/", cssHandler))
	http.Handle("/JS/", http.StripPrefix("/JS/", js_Handler))

	log.Println("El servidor se encuentra a la escucha en el puerto 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
	//Es equivalente al println solo que este terminame el programa
}