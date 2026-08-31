package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: 5 * time.Second,
}

func Projects(w http.ResponseWriter, r *http.Request) {
	resp, err := client.Get("https://api.github.com/users/CANUreal/repos")
	if err != nil {
		log.Printf("github request failed. %v", err)
		http.Error(w, "github api request failed", http.StatusBadGateway)
		// don't forget to return!!!!
		return
	}
	// don't forget about closing the body after function ends!'^!'^
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	var data any

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("%v", err)
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "	")

	if err := encoder.Encode(data); err != nil {
		log.Println(err)
	}
}
