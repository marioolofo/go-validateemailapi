package main

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	validateemail "github.com/marioolofo/go-validateemail"
)

var (
	EnvHelloHost = "VALIDATE_EMAIL_HELLO_HOST" // endereço com porta, ex: mailserver.com
	EnvEmailFrom = "VALIDATE_EMAIL_FROM" // email local, ex: exemplo@mailserver.com
	EnvApiListenIP = "VALIDATE_EMAIL_API_LISTEN_IP" // IP da API com porta, ex: localhost:80
	EnvDefaultValues = map[string]string{EnvHelloHost: "example.com", EnvEmailFrom: "test@example.com", EnvApiListenIP: ":80"}
)

// EmailValid define a resposta para a consulta de emails validos
type EmailValid struct {
	Error string `json:"error"`
	IsValid bool `json:"valid"`
}

// getEnvVar retorna o conteúdo da variável de ambiente se existir, senão retorna o defaultValue
func getEnvVar(envVar, defaultValue string) string {
	v, ok := os.LookupEnv(envVar)
	if (!ok) {
		return defaultValue
	}
	return v
}

// verifyEmail implementa o tratamento da requisição para consulta por emails
func verifyEmail(w http.ResponseWriter, r *http.Request) {
	localhost := getEnvVar(EnvHelloHost, EnvDefaultValues[EnvHelloHost])
	fromEmail := getEnvVar(EnvEmailFrom, EnvDefaultValues[EnvEmailFrom])

	vars := mux.Vars(r)
	email := vars["email"]

	ctx := validateemail.NewValidateEmail(localhost, fromEmail)
	err := ctx.Validate(email)
	if (err != nil) {
		json.NewEncoder(w).Encode(EmailValid{
			Error: err.Error(),
			IsValid: false,
		})
	} else {
		json.NewEncoder(w).Encode(EmailValid{
			Error: "",
			IsValid: true,
		})
	}
}

// setupServer configura as rotas e inicia o server
func setupServer() {
	listenIP := getEnvVar(EnvApiListenIP, EnvDefaultValues[EnvApiListenIP])
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/verify/{email}", verifyEmail).Methods("GET")
	log.Fatal(http.ListenAndServe(listenIP, myRouter))
}

// main
func main() {
	for key, _ := range EnvDefaultValues {
		_, ok := os.LookupEnv(key)
		if (!ok) {
			log.Printf(fmt.Sprintf("[WARN] %s not defined, will use default value\n", key))
		}
	}

	for key, defaultValue := range EnvDefaultValues {
		log.Printf(fmt.Sprintf("Using %s = \"%s\"", key, getEnvVar(key, defaultValue)))
	}

	setupServer()
}
