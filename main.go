package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
)

type PingResponse struct {
	Address string  `json:"address"`
	Min     float64 `json:"min"`
	Avg     float64 `json:"avg"`
	Max     float64 `json:"max"`
	Mdev    float64 `json:"mdev"`
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")

	cmd := exec.Command("ping", "-c", "4", "-i", "0.2", "-n", address)
	output, err := cmd.Output()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	statsRegex := regexp.MustCompile(`min/avg/max/mdev = ([0-9.]+)/([0-9.]+)/([0-9.]+)/([0-9.]+)`)
	stats := statsRegex.FindStringSubmatch(string(output))

	if len(stats) != 5 {
		http.Error(w, "Failed to parse ping output", http.StatusInternalServerError)
		return
	}

	min, _ := strconv.ParseFloat(stats[1], 64)
	avg, _ := strconv.ParseFloat(stats[2], 64)
	max, _ := strconv.ParseFloat(stats[3], 64)
	mdev, _ := strconv.ParseFloat(stats[4], 64)

	response := PingResponse{
		Address: address,
		Min:     min,
		Avg:     avg,
		Max:     max,
		Mdev:    mdev,
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ping", PingHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", r))
}
