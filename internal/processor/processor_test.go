package processor

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/aashish-joshi/tts-bulk/pkg/types"
)

// MockProvider is a mock TTS provider for testing.
type MockProvider struct {
	failOnGenerate bool
	callCount      int
	mu             sync.Mutex
}

func (m *MockProvider) GenerateAudio(ctx context.Context, req types.TTSRequest, outputPath string) error {
	m.mu.Lock()
	m.callCount++
	shouldFail := m.failOnGenerate
	m.mu.Unlock()

	if shouldFail {
		return os.ErrPermission
	}
	// Create an empty file to simulate audio generation
	return os.WriteFile(outputPath, []byte("mock audio data"), 0644)
}

func (m *MockProvider) Name() string {
	return "mock"
}

func (m *MockProvider) Close() error {
	return nil
}

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		cfg       Config
		wantErr   bool
		errString string
	}{
		{
			name: "valid configuration",
			cfg: Config{
				Provider:      &MockProvider{},
				OutputDir:     "test_output",
				AudioFormat:   types.FormatMP3,
				MaxConcurrent: 3,
			},
			wantErr: false,
		},
		{
			name: "missing provider",
			cfg: Config{
				Provider:      nil,
				OutputDir:     "test_output",
				AudioFormat:   types.FormatMP3,
				MaxConcurrent: 3,
			},
			wantErr:   true,
			errString: "provider is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up
			defer os.RemoveAll(tt.cfg.OutputDir)

			proc, err := New(tt.cfg)

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

			if proc == nil {
				t.Errorf("New() returned nil processor")
			}

			proc.Close()
		})
	}
}

func TestProcessCSV(t *testing.T) {
	// Create a temporary CSV file
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")
	csvContent := `"label1","Hello World"
"label2","Test message"`

	if err := os.WriteFile(csvPath, []byte(csvContent), 0644); err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	outputDir := filepath.Join(tmpDir, "audio")

	tests := []struct {
		name           string
		provider       *MockProvider
		csvPath        string
		wantErr        bool
		wantSuccessful int
		wantFailed     int
	}{
		{
			name:           "successful processing",
			provider:       &MockProvider{failOnGenerate: false},
			csvPath:        csvPath,
			wantErr:        false,
			wantSuccessful: 2,
			wantFailed:     0,
		},
		{
			name:           "file not found",
			provider:       &MockProvider{failOnGenerate: false},
			csvPath:        filepath.Join(tmpDir, "nonexistent.csv"),
			wantErr:        true,
			wantSuccessful: 0,
			wantFailed:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proc, err := New(Config{
				Provider:      tt.provider,
				OutputDir:     outputDir,
				AudioFormat:   types.FormatMP3,
				MaxConcurrent: 2,
			})
			if err != nil {
				t.Fatalf("Failed to create processor: %v", err)
			}
			defer proc.Close()

			ctx := context.Background()
			results, err := proc.ProcessCSV(ctx, tt.csvPath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ProcessCSV() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ProcessCSV() unexpected error: %v", err)
				return
			}

			successful := 0
			failed := 0
			for _, result := range results {
				if result.Error == nil {
					successful++
				} else {
					failed++
				}
			}

			if successful != tt.wantSuccessful {
				t.Errorf("ProcessCSV() successful = %d, want %d", successful, tt.wantSuccessful)
			}

			if failed != tt.wantFailed {
				t.Errorf("ProcessCSV() failed = %d, want %d", failed, tt.wantFailed)
			}
		})
	}
}

func TestReadCSV(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name       string
		csvContent string
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "valid CSV",
			csvContent: "\"label1\",\"text1\"\n\"label2\",\"text2\"",
			wantCount:  2,
			wantErr:    false,
		},
		{
			name:       "CSV with invalid rows",
			csvContent: "\"label1\",\"text1\"\n\"invalid\"\n\"label2\",\"text2\"",
			wantCount:  0,
			wantErr:    true, // CSV parser will fail on invalid row
		},
		{
			name:       "empty CSV",
			csvContent: "",
			wantCount:  0,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			csvPath := filepath.Join(tmpDir, "test_"+tt.name+".csv")
			if err := os.WriteFile(csvPath, []byte(tt.csvContent), 0644); err != nil {
				t.Fatalf("Failed to create test CSV: %v", err)
			}

			proc := &Processor{}
			records, err := proc.readCSV(csvPath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("readCSV() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("readCSV() unexpected error: %v", err)
				return
			}

			if len(records) != tt.wantCount {
				t.Errorf("readCSV() count = %d, want %d", len(records), tt.wantCount)
			}
		})
	}
}
