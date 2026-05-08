package builtin

import (
	"context"
	"testing"
)

func TestCalculatorInvoke(t *testing.T) {
	c := Calculator{}
	tests := []struct {
		args string
		want string
	}{
		{`{"a":1,"b":2,"op":"add"}`, "3"},
		{`{"a":10,"b":4,"op":"sub"}`, "6"},
		{`{"a":3,"b":4,"op":"mul"}`, "12"},
		{`{"a":8,"b":2,"op":"div"}`, "4"},
	}
	for _, tt := range tests {
		got, err := c.Invoke(context.Background(), tt.args)
		if err != nil {
			t.Fatalf("args %s: %v", tt.args, err)
		}
		if got != tt.want {
			t.Fatalf("args %s: got %q want %q", tt.args, got, tt.want)
		}
	}
}

func TestCalculatorDivideByZero(t *testing.T) {
	c := Calculator{}
	_, err := c.Invoke(context.Background(), `{"a":1,"b":0,"op":"div"}`)
	if err == nil {
		t.Fatal("expected error for division by zero")
	}
}
