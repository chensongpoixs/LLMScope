package handler

import (
	"testing"
)

func TestExtractParamsFromModelName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Gemma 3B", "gemma-4-3b-it-Q4_0.gguf", "3B"},
		{"LLaMA 7B", "llama-2-7b-chat.Q4_K_M.gguf", "7B"},
		{"LLaMA 13B", "llama-2-13b.Q5_K_S.gguf", "13B"},
		{"Mistral 7B", "mistral-7b-instruct-v0.2.Q4_0.gguf", "7B"},
		{"Qwen 1.5B", "qwen-1.5b-chat.Q4_K_M.gguf", "1.5B"},
		{"Unknown", "model-unknown.gguf", "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractParamsFromModelName(tt.input)
			if result != tt.expected {
				t.Errorf("extractParamsFromModelName(%s) = %s; want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractQuantizationFromModelName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Q4_0", "model.Q4_0.gguf", "Q4_0"},
		{"Q4_K_M", "model.Q4_K_M.gguf", "Q4_K_M"},
		{"Q5_K_S", "model.Q5_K_S.gguf", "Q5_K_S"},
		{"F16", "model.F16.gguf", "F16"},
		{"Unknown", "model.gguf", "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractQuantizationFromModelName(tt.input)
			if result != tt.expected {
				t.Errorf("extractQuantizationFromModelName(%s) = %s; want %s", tt.input, result, tt.expected)
			}
		})
	}
}
