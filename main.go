package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	trek "real/assets/struct"
	utils "real/assets/utils"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var db *sql.DB

type Denomination struct {
	FormName string `json:"formName"`
	Username string `json:"Username"`
	Other    string `json:"Other"`
	Content  string `json:"Content"`
}
type FormDataRegister struct {
	FormName        string `json:"formName"`
	Username        string `json:"username"`
	Email           string `json:"mail"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type FormDataLogin struct {
	FormName string `json:"formName"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type FormDataPost struct {
	FormName string `json:"formName"`
	Username string `json:"Username"`
	Subject  string `json:"Subject"`
	Content  string `json:"Content"`
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./db/forum.sqlite")
	if err != nil {
		log.Fatal(err)
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	servePage(w, r, "templates/index.html")
}

func servePage(w http.ResponseWriter, r *http.Request, pageName string) {
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pagePath := filepath.Join(wd, pageName)
	http.ServeFile(w, r, pagePath)
}

func loginHandler(conn *websocket.Conn, message []byte) {

	var formData FormDataLogin
	err := json.Unmarshal(message, &formData)
	if err != nil {
		fmt.Println("Erreur lors de l'analyse des données JSON:", err)
		return
	}

	// Gérez les données du formulaire de connexion
	fmt.Println("Données du formulaire de connexion:")
	fmt.Println("Nom d'utilisateur:", formData.Username)
	fmt.Println("Mot de passe:", formData.Password)

	user, err := utils.GetUserByUsername(db, formData.Username)
	if err != nil {
		fmt.Printf("Can't find %s in the database.", formData.Username)
	}

	if user != nil && user.Password == formData.Password {
		fmt.Printf("%s was successfully logged.\n", user.Username)
		utils.Activation(db, formData.Username)

		response := trek.LoginResponse{Success: true, Message: "Everything is fine.", Name: "Login"}
		responseData, err := json.Marshal(response)
		fmt.Println(string(responseData))
		if err != nil {
			fmt.Println("Problème pour tout remettre en JSON.")
			return
		}
		err = conn.WriteMessage(websocket.TextMessage, responseData)
		if err != nil {
			fmt.Println("Problème pour l'envoi du JSON vers le script.")
			return
		}
	} else {
		fmt.Println("Impossible de se connecter avec ce pseudo.")
		response := trek.LoginResponse{Success: false, Message: "Nom ou MDP invalide.", Name: "Login"}
		responseData, err := json.Marshal(response)
		if err != nil {
			fmt.Println("Erreur pour recompiler le message d'échec en JSON.")
			return
		}
		err = conn.WriteMessage(websocket.TextMessage, responseData)
		if err != nil {
			fmt.Println("Erreur lors du renvoi du JSON d'erreur dans le script.")
		}
	}
}

var activeConn int

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade de la connexion HTTP vers une connexion WebSocket
	activeConn++
	utils.Deactivation(db)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Erreur lors de l'upgrade de la connexion WebSocket:", err)
		return
	}
	defer func() {
		// Decrement the activeConnections counter when the connection is closed
		activeConn--
		// Check if there are no more active connections
		if activeConn == 0 {
			// Perform deactivation logic when there are no active connections
			err := utils.Deactivation(db)
			if err != nil {
				fmt.Println("Error deactivating all users:", err)
			}
		}
		conn.Close()
	}()
	// fmt.Println("je rentre2")
	// Boucle pour lire les messages WebSocket
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Erreur lors de la lecture du message WebSocket:", err)
			break
		}
		var nomForm Denomination
		var postForm FormDataPost
		err = json.Unmarshal(message, &nomForm)
		if err != nil {
			fmt.Println("Erreur lors de l'analyse des données JSON:", err)
			continue
		}
		// Traitez les données en fonction du nom du formulaire
		fmt.Println(activeConn)
		// fmt.Println("form", nomForm.FormName, nomForm.Username)
		switch nomForm.FormName {
		case "register":
			registerHandler(message)
		case "login":
			loginHandler(conn, message)
		case "posts":
			sendPostsToClients(conn)
		case "reset":
			utils.Reboot(db)
		case "usershunt":
			friends, err := utils.GetFriends(db, nomForm.Username)
			if err != nil {
				fmt.Println("Error getting friends:", err)
				return
			}
			messageData := trek.UsersArray{Name: "Friends", Users: friends}
			message, _ := json.Marshal(messageData)
			err = conn.WriteMessage(websocket.TextMessage, message)
			// fmt.Println(friends)
			if err != nil {
				fmt.Println("Error sending friends list:", err)
				return
			}
		case "userStatus":
			update, err := utils.GetStatus(db, nomForm.Username)
			if err != nil {
				fmt.Println("Error updating status:", err)
				return
			}
			messageData := trek.NewStatus{Name: "userStatus", Checks: update}
			message, _ := json.Marshal(messageData)
			err = conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				fmt.Println("Error pour envoyer la liste d'amis:", err)
				return
			}
		case "discussions":
			fmt.Println(nomForm.Other)
			discussionHandler(conn, nomForm)
		case "chatEnvoy":
			fmt.Println("Etape 2")
			fmt.Println(nomForm)
			utils.NewMessage(db, nomForm.Username, nomForm.Other, nomForm.Content)
		case "delog":
			utils.Delog(db, nomForm.Username)
		case "postCreation":
			json.Unmarshal(message, &postForm)
			utils.CreatePost(db, postForm.Username, postForm.Subject, postForm.Content)
		default:
			fmt.Println("Nom de formulaire non reconnu:", nomForm.FormName)
		}
	}
}

func discussionHandler(conn *websocket.Conn, data Denomination) {
	convs, err := utils.GetDiscussion(db, data.Username, data.Other)
	if err != nil {
		fmt.Println("Erreur pour chopper les convs:", err)
		return
	}
	messageData := trek.MessageOuter{Name: "chatHistory", Chats: convs}
	message, _ := json.Marshal(messageData)
	err = conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		fmt.Println("Erreur pour envoyer l'historique des chats:", err)
		return
	}
}

func registerHandler(message []byte) {
	var formData FormDataRegister
	err := json.Unmarshal(message, &formData)
	if err != nil {
		fmt.Println("Erreur lors de l'analyse des données JSON:", err)
		return
	}
	utils.InsertUser(db, formData.Username, formData.Email, formData.Password, formData.ConfirmPassword)
}

func main() {
	initDB()
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/ws", wsHandler)
	fmt.Println("Serveur WebSocket démarré sur ws://localhost:8080/ws")
	fmt.Println("Vrai serveur démarré sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
	defer db.Close()
}

func sendPostsToClients(conn *websocket.Conn) {
	posts, err := utils.GetPosts(db)
	if err != nil {
		fmt.Println("Erreur pour chopper les données de GetPosts()!")
		return
	}
	// fmt.Println(posts)
	message := trek.PostArray{Name: "Posts", Posts: posts}
	postData, _ := json.Marshal(message)
	err = conn.WriteMessage(websocket.TextMessage, postData)
	if err != nil {
		fmt.Println("Erreur pour renvoyer les données vers le JS!")
		return
	}
}
