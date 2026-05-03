package chat

import "testing"

func TestGetOptions(t *testing.T) {
	t.Run("nil base returns zero value options", func(t *testing.T) {
		opts := GetOptions(nil)
		if opts == nil {
			t.Fatal("expected non-nil options")
		}
		if opts.Temperature != 0 {
			t.Fatalf("expected Temperature 0, got %v", opts.Temperature)
		}
		if opts.MaxTokens != 0 {
			t.Fatalf("expected MaxTokens 0, got %v", opts.MaxTokens)
		}
		if opts.MaxCompletionTokens != 0 {
			t.Fatalf("expected MaxCompletionTokens 0, got %v", opts.MaxCompletionTokens)
		}
		if opts.Model != "" {
			t.Fatalf("expected Model empty, got %q", opts.Model)
		}
	})

	t.Run("base preserved when no options provided", func(t *testing.T) {
		base := &Options{Temperature: 0.8, MaxTokens: 256, MaxCompletionTokens: 32, Model: "gpt-test"}
		opts := GetOptions(base)
		if opts != base {
			t.Fatal("expected returned options to be the same pointer as base")
		}
		if opts.Temperature != 0.8 || opts.MaxTokens != 256 || opts.MaxCompletionTokens != 32 || opts.Model != "gpt-test" {
			t.Fatal("expected base values to remain unchanged")
		}
	})

	t.Run("ignore option with nil apply", func(t *testing.T) {
		base := &Options{Temperature: 0.9, MaxTokens: 512, MaxCompletionTokens: 64, Model: "gpt-nil"}
		opts := GetOptions(base, Option{})
		if opts.Temperature != 0.9 || opts.MaxTokens != 512 || opts.MaxCompletionTokens != 64 || opts.Model != "gpt-nil" {
			t.Fatal("expected nil apply option to be ignored")
		}
	})

	t.Run("zero values are applied correctly", func(t *testing.T) {
		base := &Options{Temperature: 0.5, MaxTokens: 1024, MaxCompletionTokens: 128, Model: "gpt-nonzero"}
		zeroOpt := Option{apply: func(opts *Options) {
			opts.Temperature = 0
			opts.MaxTokens = 0
			opts.MaxCompletionTokens = 0
			opts.Model = ""
		}}
		opts := GetOptions(base, zeroOpt)
		if opts.Temperature != 0 {
			t.Fatalf("expected Temperature 0, got %v", opts.Temperature)
		}
		if opts.MaxTokens != 0 {
			t.Fatalf("expected MaxTokens 0, got %v", opts.MaxTokens)
		}
		if opts.MaxCompletionTokens != 0 {
			t.Fatalf("expected MaxCompletionTokens 0, got %v", opts.MaxCompletionTokens)
		}
		if opts.Model != "" {
			t.Fatalf("expected Model empty, got %q", opts.Model)
		}
	})

	t.Run("multiple options apply in order and override values", func(t *testing.T) {
		base := &Options{Temperature: 0.2, MaxTokens: 128, MaxCompletionTokens: 16, Model: "initial"}
		optA := Option{apply: func(opts *Options) {
			opts.Temperature = 0.7
			opts.Model = "step-a"
		}}
		optB := Option{apply: func(opts *Options) {
			opts.MaxTokens = 256
			opts.MaxCompletionTokens = 64
			opts.Model = "step-b"
		}}
		opts := GetOptions(base, optA, optB)
		if opts.Temperature != 0.7 {
			t.Fatalf("expected Temperature 0.7, got %v", opts.Temperature)
		}
		if opts.MaxTokens != 256 {
			t.Fatalf("expected MaxTokens 256, got %v", opts.MaxTokens)
		}
		if opts.MaxCompletionTokens != 64 {
			t.Fatalf("expected MaxCompletionTokens 64, got %v", opts.MaxCompletionTokens)
		}
		if opts.Model != "step-b" {
			t.Fatalf("expected Model step-b, got %q", opts.Model)
		}
	})
}
