package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/rest"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/speak"
)

func checkDeepgramKey() error {
	// Check if the Deepgram API key is set
	if os.Getenv("DEEPGRAM_API_KEY") == "" {
		return fmt.Errorf("DEEPGRAM_API_KEY environment variable is not set")
	}
	return nil
}

func main() {

	var dgContainer, dgEncoding string

	// Define the commandline flags
	dgModelName := flag.String("model", "aura-asteria-en", "Deepgram model name. Defaults to aura-asteria-en")
	dgFileFormat := flag.String("format", "mp3", "File format for the generated audio files. Defaults to mp3")
	outputFolder := flag.String("output", "audio", "Output folder for the generated audio files.")
	csvLocation := flag.String("csv", "scripts.csv", "Location of the CSV file containing the scripts to convert to audio.")
	// parse the commandline flags
	flag.Parse()

	// If the DG file format is neither wav nor mp3, raise error.
	if strings.ToLower(*dgFileFormat) != "wav" && strings.ToLower(*dgFileFormat) != "mp3" {
		fmt.Println("Invalid file format. Only wav and mp3 are supported.")
		return
	}

	if strings.ToLower(*dgFileFormat) == "wav" {
		dgContainer = "wav"
		dgEncoding = "linear16"
	} else {
		dgContainer = ""
		dgEncoding = "mp3"
	}
	// First check if the Deepgram API key is set
	if dgErr := checkDeepgramKey(); dgErr != nil {
		fmt.Println(dgErr)
		return
	}

	client.Init(client.InitLib{
		LogLevel: client.LogLevelErrorOnly,
	})

	ctx := context.Background()

	// set the Transcription options
	options := &interfaces.SpeakOptions{
		Model:     strings.ToLower(*dgModelName),
		Container: dgContainer,
		Encoding:  dgEncoding,
	}

	// create a Deepgram client
	c := client.NewRESTWithDefaults()
	dg := api.New(c)

	// Check if the file exists
	if _, err := os.Stat(*csvLocation); os.IsNotExist(err) {
		fmt.Print("The file does not exist. Please enter the correct file path.")
		return
	}

	// Open the CSV file
	file, err := os.Open(*csvLocation)
	if err != nil {
		fmt.Printf("Error opening CSV file: %s\n", err)
		return
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading CSV file: %s\n", err)
		return
	}

	// Create the "audio" directory if it doesn't exist
	audioDir := strings.ToLower(*outputFolder)
	if _, err := os.Stat(audioDir); os.IsNotExist(err) {
		err := os.Mkdir(audioDir, 0755)
		if err != nil {
			fmt.Printf("Could not create directory: %s\n", err)
			return
		}
	}

	// Limit the number of concurrently running goroutines
	const maxGoroutines = 3
	guard := make(chan struct{}, maxGoroutines)
	var wg sync.WaitGroup

	for i, record := range records {
		if len(record) < 2 {
			fmt.Println("Skipping invalid record:", record)
			continue
		}
		label := record[0]
		script := record[1]

		// Acquire a slot in the goroutine pool
		guard <- struct{}{}
		wg.Add(1)

		go func(i int, label, script string) {
			defer wg.Done()

			// Perform TTS and save to disk
			audioPath := filepath.Join(audioDir, fmt.Sprintf("%s.%s", label, *dgFileFormat))
			err := generateTTSAndSave(ctx, dg, script, options, audioPath)
			if err != nil {
				fmt.Printf("Could not generate TTS for row %v - %v: %v\n", i, label, err)
			} else {
				fmt.Printf("TTS generated for %v\n", label)
			}

			// Release a slot in the goroutine pool
			<-guard
		}(i, label, script)
	}

	// Wait for all goroutines to complete
	wg.Wait()
}

func generateTTSAndSave(ctx context.Context, dg *api.Client, script string, options *interfaces.SpeakOptions, audioPath string) error {
	// Generate TTS data
	_, err := dg.ToSave(ctx, audioPath, script, options)
	if err != nil {
		return err
	}

	// Wait for 1 second before returning
	// This is to avoid rate limiting
	time.Sleep(time.Second)

	return nil
}
