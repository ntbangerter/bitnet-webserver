package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
)

type Request struct {
	Prompt string `json:"input"`
}

type Response struct {
	Output string `json:"output"`
}

func prompt(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := exec.Command(
		"/home/tanner/BitNet/build/bin/llama-cli",
		"-m", "/home/tanner/BitNet/models/bitnet_b1_58-large/ggml-model-i2_s.gguf",
		"-p", req.Prompt,
		"-n", "32")

	output, err := cmd.Output()

	if err != nil {
		log.Fatalf("Command execution failed: %s\nOutput: %s", err, output)
	}

	res := Response{Output: string(output)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	http.HandleFunc("/prompt", prompt)
	http.HandleFunc("/", serveIndex)

	http.ListenAndServe("0.0.0.0:5000", nil)
}
