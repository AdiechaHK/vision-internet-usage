package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello word !")
}

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", hello)
	fmt.Println("Listening on port: " + port)
	http.ListenAndServe(":"+port, nil)
}
