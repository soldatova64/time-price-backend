package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", pageHome)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func pageHome(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Привет!")
}
