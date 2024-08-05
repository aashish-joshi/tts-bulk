# Bulk TTS Generator

This project uses Deepgram's TTS service to generate audio files in bulk from a CSV.

The project is written in Go and executables are provided for Windows, MacOS, and Linux.

## Usage

1. Get an API key from [Deepgram](https://www.deepgram.com/) and add it to your environment variables as `DEEPGRAM_API_KEY`.
2. Create a CSV file with the following columns (use the provided `sample-scripts.csv` as a template)
    - `label`: The label for the script. This will be used as the file name.
    - `script`: The text to be converted to speech.
3. Download the executable for your OS from the [releases](https://github.com/aashish-joshi/tts-bulk/releases) page.
4. The tool will try and read the csv locally from `scripts.csv`. If it doesn't exist, it will ask for the path to the CSV file.

### Commandline flags

The following flags are supported at the moment.

- `-format`: The format of the audio file. Supported formats are `wav` and `mp3`. Default is `mp3`.
- `-output`: The output directory where the audio files will be saved. Default is `audio/`.
- `-csv`: The path to the CSV file. Default is `scripts.csv`.
- `-output`: The output directory where the audio files will be saved. Default is `audio/`.

### Example

1. Generate mp3 files in the default location.
    
    ```bash
    ./tts-bulk
    ```

2. Generate wav files in a custom location.

    ```bash
    ./tts-bulk -format=wav -output=/path/to/output
    ```
