package auth

import (
	"net/http"
)

// import "github.com/go-chi/chi/v5"

func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
	// regReq := models.RegisterRequest{}
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
