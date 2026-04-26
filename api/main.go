package main

import (
	"api/handler"
	"api/middleware"
	"api/playground"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// route handlers
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, Go!")
}

func timeHandler(w http.ResponseWriter, _ *http.Request) {
	now := time.Now()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonStr := fmt.Sprintf(`{"id":"%s"}`, now.String())
	w.Write([]byte(jsonStr))
}

func echoHanlder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonStr := fmt.Sprintf(`{"data":"%s"}`, body)
	w.Write([]byte(jsonStr))
}

func main() {

	// Route
	mux := http.NewServeMux()

	// Register a new route - handleFunc
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("GET /time", timeHandler)
	mux.HandleFunc("POST /echo", echoHanlder)
	mux.HandleFunc("GET /users", handler.ListUsers)
	mux.HandleFunc("GET /users/{id}", handler.GetUser)
	mux.Handle("POST /users", middleware.Chain(http.HandlerFunc(handler.ListUsers), middleware.Auth))
	mux.HandleFunc("DELETE /users/{id}", handler.DeleteUser)
	mux.HandleFunc("PUT /users/{id}", handler.UpdateUser)
	mux.HandleFunc("GET /playground", playground.Handler)

	// Inline route handler
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	stack := middleware.Chain(mux, middleware.Logger, middleware.Timeout(8*time.Second))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      stack,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("Server is listening on :8080")
	log.Fatal(server.ListenAndServe())

}
