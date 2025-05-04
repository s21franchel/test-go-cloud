package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Response from backend %s", port)
	})

	fmt.Printf("Backend server started on :%s\n", port)
	http.ListenAndServe(":"+port, nil)
}
