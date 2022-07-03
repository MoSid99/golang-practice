package main

import (
	"example/storybook"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func check_err(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	port := flag.Int("port", 3000, "Server port")
	// Get the file we want to read
	filepath := flag.String("filepath", "gopher.json", "file path")
	flag.Parse()
	fmt.Printf("Reading: %s", *filepath)

	// Open the file
	file, err := os.Open(*filepath)
	check_err(err)

	// JSON decoding
	story, err := storybook.JSONStory(file)
	check_err(err)

	handler := storybook.StoryHandler(story)
	fmt.Printf("Server starting on %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handler))
}
