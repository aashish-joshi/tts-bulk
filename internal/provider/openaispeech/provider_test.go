package openaispeech

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aashish-joshi/tts-bulk/pkg/types"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid config with all fields",
			cfg: Config{
				BaseURL:     "http://localhost:8000",
				Model:       "tts-1",
				Voice:       "alloy",
				Format:      "mp3",
				RateLimitMs: 0,
			},
			wantErr: false,
		},
		{
			name:    "defaults applied when fields empty",
			cfg:     Config{},
			wantErr: false,
		},
		{
			name: "trailing slash stripped from BaseURL",
			cfg: Config{
				BaseURL: "http://localhost:8000/",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := New(tt.cfg)
			if tt.wantErr {
				if err == nil {
					t.Errorf("New() expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("New() unexpected error: %v", err)
				return
			}
			if p == nil {
				t.Errorf("New() returned nil provider")
			}
		})
	}
}

func TestName(t *testing.T) {
	p, err := New(Config{})
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
	if p.Name() != "openaispeech" {
		t.Errorf("Name() = %q, want %q", p.Name(), "openaispeech")
	}
}

func TestClose(t *testing.T) {
	p, err := New(Config{})
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
	if err := p.Close(); err != nil {
		t.Errorf("Close() unexpected error: %v", err)
	}
}

func TestGenerateAudio(t *testing.T) {
	tests := []struct {
		name           string
		serverStatus   int
		serverBody     []byte
		req            types.TTSRequest
		apiKey         string
		wantErr        bool
		wantErrContain string
	}{
		{
			name:         "success mp3",
			serverStatus: http.StatusOK,
			serverBody:   []byte("fake-mp3-audio-data"),
			req:          types.TTSRequest{Text: "Hello world", Label: "hello"},
			wantErr:      false,
		},
		{
			name:         "success wav",
			serverStatus: http.StatusOK,
			serverBody:   []byte("fake-wav-audio-data"),
			req:          types.TTSRequest{Text: "Test text", Label: "test"},
			wantErr:      false,
		},
		{
			name:           "empty text returns error",
			serverStatus:   http.StatusOK,
			serverBody:     []byte{},
			req:            types.TTSRequest{Text: "", Label: "empty"},
			wantErr:        true,
			wantErrContain: "text cannot be empty",
		},
		{
			name:           "server returns HTTP error",
			serverStatus:   http.StatusInternalServerError,
			serverBody:     []byte("internal server error"),
			req:            types.TTSRequest{Text: "Hello", Label: "err"},
			wantErr:        true,
			wantErrContain: "500",
		},
		{
			name:         "request includes Authorization header when API key set",
			serverStatus: http.StatusOK,
			serverBody:   []byte("audio-data"),
			req:          types.TTSRequest{Text: "Secure request", Label: "secure"},
			apiKey:       "test-api-key",
			wantErr:      false,
		},
		{
			name:         "no API key works for local server",
			serverStatus: http.StatusOK,
			serverBody:   []byte("audio-data"),
			req:          types.TTSRequest{Text: "Local request", Label: "local"},
			apiKey:       "",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedReq *http.Request
			var capturedBody speechRequest

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedReq = r
				if err := json.NewDecoder(r.Body).Decode(&capturedBody); err != nil {
					http.Error(w, "bad request", http.StatusBadRequest)
					return
				}
				w.WriteHeader(tt.serverStatus)
				w.Write(tt.serverBody) //nolint:errcheck
			}))
			defer srv.Close()

			p, err := New(Config{
				BaseURL:     srv.URL,
				Model:       "tts-1",
				Voice:       "alloy",
				Format:      "mp3",
				APIKey:      tt.apiKey,
				RateLimitMs: 0,
			})
			if err != nil {
				t.Fatalf("New() unexpected error: %v", err)
			}

			outputPath := filepath.Join(t.TempDir(), tt.req.Label+".mp3")

			err = p.GenerateAudio(context.Background(), tt.req, outputPath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GenerateAudio() expected error but got none")
				}
				if tt.wantErrContain != "" && err != nil {
					if !strings.Contains(err.Error(), tt.wantErrContain) {
						t.Errorf("GenerateAudio() error = %q, want to contain %q", err.Error(), tt.wantErrContain)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("GenerateAudio() unexpected error: %v", err)
				return
			}

			// Verify the file was written with correct content
			data, err := os.ReadFile(outputPath)
			if err != nil {
				t.Errorf("failed to read output file: %v", err)
				return
			}
			if string(data) != string(tt.serverBody) {
				t.Errorf("file content = %q, want %q", string(data), string(tt.serverBody))
			}

			// Verify request headers
			if capturedReq != nil {
				if capturedReq.Header.Get("Content-Type") != "application/json" {
					t.Errorf("Content-Type = %q, want %q", capturedReq.Header.Get("Content-Type"), "application/json")
				}
				if tt.apiKey != "" {
					wantAuth := "Bearer " + tt.apiKey
					if capturedReq.Header.Get("Authorization") != wantAuth {
						t.Errorf("Authorization = %q, want %q", capturedReq.Header.Get("Authorization"), wantAuth)
					}
				} else {
					if capturedReq.Header.Get("Authorization") != "" {
						t.Errorf("Authorization header should not be set when no API key")
					}
				}

				// Verify request body fields
				if capturedBody.Input != tt.req.Text {
					t.Errorf("request body input = %q, want %q", capturedBody.Input, tt.req.Text)
				}
				if capturedBody.Model != "tts-1" {
					t.Errorf("request body model = %q, want %q", capturedBody.Model, "tts-1")
				}
				if capturedBody.Voice != "alloy" {
					t.Errorf("request body voice = %q, want %q", capturedBody.Voice, "alloy")
				}
			}
		})
	}
}

func TestGenerateAudioEmptyOutputPath(t *testing.T) {
	p, err := New(Config{BaseURL: "http://localhost:8000"})
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	err = p.GenerateAudio(context.Background(), types.TTSRequest{Text: "hello"}, "")
	if err == nil {
		t.Errorf("GenerateAudio() expected error for empty output path")
	}
}
