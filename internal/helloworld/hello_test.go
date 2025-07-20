package helloworld

import "testing"

func TestGreet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"default", "", "Hello, World!"},
		{"name", "Go", "Hello, Go!"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := Greet(tt.input)
			if got != tt.want {
				t.Fatalf("Greet(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
