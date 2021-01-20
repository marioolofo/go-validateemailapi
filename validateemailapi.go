package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	validateemail "github.com/marioolofo/go-validateemail"
)

type EmailValid struct {
	Error string `json:"error"`
	IsValid bool `json:"valid"`
}

func verifyEmail(w http.ResponseWriter, r *http.Request) {
	localhost, ok := os.LookupEnv("VALIDATE_EMAIL_LOCALHOST")
	if (!ok) {
		json.NewEncoder(w).Encode(EmailValid{
			Error: "Env var VALIDATE_EMAIL_LOCALHOST not configured",
			IsValid: false,
		})
		return
	}
	fromEmail, ok := os.LookupEnv("VALIDATE_EMAIL_FROM")
	if (!ok) {
		json.NewEncoder(w).Encode(EmailValid{
			Error: "Env var VALIDATE_EMAIL_FROM not configured",
			IsValid: false,
		})
		return
	}

	vars := mux.Vars(r)
	email := vars["email"]

	ctx := validateemail.NewValidateEmail(localhost, fromEmail)
	result := ctx.Validate(email)
	if (result != nil) {
		json.NewEncoder(w).Encode(EmailValid{
			Error: result.Error(),
			IsValid: false,
		})
	} else {
		json.NewEncoder(w).Encode(EmailValid{
			Error: "",
			IsValid: true,
		})
	}
}

func setupServer() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/verify/{email}", verifyEmail).Methods("GET")
	log.Fatal(http.ListenAndServe(":80", myRouter))
}

func checkEnv() {
	_, ok := os.LookupEnv("VALIDATE_EMAIL_LOCALHOST")
	if (!ok) {
		log.Fatal("VALIDATE_EMAIL_LOCALHOST not defined")
	}
	_, ok = os.LookupEnv("VALIDATE_EMAIL_FROM")
	if (!ok) {
		log.Fatal("VALIDATE_EMAIL_FROM not defined")
	}
}

func main() {
	checkEnv()
	setupServer()
}
