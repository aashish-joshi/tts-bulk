# Bulk TTS Generator

This project uses Deepgram's TTS service to generate audio files in bulk from a CSV.

The project is written in Go and executables are provided for Windows, MacOS, and Linux.

## Usage

1. Get an API key from [Deepgram](https://www.deepgram.com/) and add it to your environment variables as `DEEPGRAM_API_KEY`.
2. Create a CSV file with the following columns:
    - `label`: The label for the script. This will be used as the file name.
    - `script`: The text to be converted to speech.
3. Run the executable. It will try and read the csv locally from `./scripts.csv`. If it doesn't exist, it will ask for the path to the CSV file.
4. The audio files will be saved in the `audio/` directory.

