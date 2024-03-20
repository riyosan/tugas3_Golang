package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

func updateStatus() {
	for {
		status := Status{
			Water: rand.Intn(100) + 1,
			Wind:  rand.Intn(100) + 1,
		}

		data, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			continue
		}

		err = ioutil.WriteFile("status.json", data, 0644)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			continue
		}

		time.Sleep(15 * time.Second)
	}
}

func main() {
	go updateStatus()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile("status.json")
		if err != nil {
			http.Error(w, "Error reading status", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
