package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type ArchiveRequestPayload struct {
	TimeZone string `json:"timeZone"`
	Command  string `json:"command"`
	Schedule string `json:"schedule"`
}

func buildArchiveRequestHandler(writer http.ResponseWriter, req *http.Request) {
	// Validate payload arg
	payloadArgs := req.URL.Query()["payload"]
	if payloadArgs == nil || len(payloadArgs) == 0 {
		writer.Write([]byte("error no payload"))
		return
	}

	// Deserialize json
	payloadString := payloadArgs[0]
	payload := ArchiveRequestPayload{}
	err := json.Unmarshal([]byte(payloadString), &payload)
	if err != nil {
		writer.Write([]byte(fmt.Sprintf("error invalid payload: %s", err)))
		return
	}

	archiveName, err := createArchive(payload)
	if err != nil {
		writer.Write([]byte(fmt.Sprintf("error creating archive: %s", err)))
	}

	file, err := os.Open(archiveName)
	if err != nil {
		writer.Write([]byte(fmt.Sprintf("error opening archive to send: %s", err)))
	}
	defer file.Close()
	defer os.Remove(archiveName)

	io.Copy(writer, file)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("static/")))
	http.HandleFunc("/generate", buildArchiveRequestHandler)
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
