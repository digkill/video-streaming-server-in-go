package main

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", streamHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	//http.ServeFile(w, r, "video.mp4")
	file, header, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Failed to read video file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	partReader := multipart.NewReader(file, strconv.FormatInt(header.Size, 10))
	part, err := partReader.NextPart()
	if err != nil {
		http.Error(w, "Failed to read video part", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", part.Header.Get("Content-Type"))
	w.Header().Set("Content-Disposition", "inline")
	w.WriteHeader(http.StatusOK)

	// Start streaming the video
	_, err = io.Copy(w, part)
	if err != nil {
		log.Println("Failed to stream video:", err)
	}
}
