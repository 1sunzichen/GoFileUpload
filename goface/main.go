package main

import (
	"github.com/Kagami/go-face"
	"log"
)
const dataDir = "testdata"
func main(){
	rec, err := face.NewRecognizer(dataDir)
	if err != nil {
		log.Fatalf("Can't init face recognizer: %v", err)
	}
	// Free the resources when you're finished.
	defer rec.Close()
}
