package version

import "testing"

func TestPrint(t *testing.T) {
	t.Parallel()
	Version = "1.0.0"
	Commit = "abc123"
	BuildDate = "2025-01-01"

	want := "1.0.0 (abc123 2025-01-01)"
	got := Print()
	if got != want {
		t.Fatalf("Print() = %q, want %q", got, want)
	}
}
