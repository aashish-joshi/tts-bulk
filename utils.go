package main

import (
	"fmt"
	"os"
	"strings"
)

func checkDeepgramKey() error {
	// Check if the Deepgram API key is set
	if os.Getenv("DEEPGRAM_API_KEY") == "" {
		return fmt.Errorf("DEEPGRAM_API_KEY environment variable is not set")
	}
	return nil
}

func setupAudioDir(dirName string) (string, error) {
	// Create the "audio" directory if it doesn't exist
	audioDir := strings.ToLower(dirName)
	if _, err := os.Stat(audioDir); os.IsNotExist(err) {
		err := os.Mkdir(audioDir, 0755)
		if err != nil {
			return "", err
		}
	}
	return audioDir, nil
}
