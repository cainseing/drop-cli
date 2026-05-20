package display

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"
)

func TestUcFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"a", "A"},
		{"hello", "Hello"},
		{"Hello", "Hello"},
		{"123", "123"},
		{"ü", "Ü"},
	}

	for _, test := range tests {
		result := ucFirst(test.input)
		if result != test.expected {
			t.Errorf("ucFirst(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}

func TestPrintError(t *testing.T) {
	tests := []struct {
		message  string
		err      error
		contains []string
	}{
		{"test message", errors.New("test error"), []string{"ERROR", "Test message (test error)"}},
		{"test message", nil, []string{"ERROR", "Test message"}},
		{"", errors.New("test error"), []string{"ERROR", "Test error"}},
		{"", nil, []string{"ERROR"}}, // Should not print anything extra
	}

	for _, test := range tests {
		// Capture stderr
		oldStderr := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w

		PrintError(test.message, test.err)

		w.Close()
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)
		output := buf.String()
		os.Stderr = oldStderr

		for _, substr := range test.contains {
			if !strings.Contains(output, substr) {
				t.Errorf("PrintError(%q, %v) output should contain %q, got: %q", test.message, test.err, substr, output)
			}
		}
	}
}

func TestPrintProperty(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintProperty("LABEL", "value")

	w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()
	os.Stdout = oldStdout

	if !strings.Contains(output, "LABEL") || !strings.Contains(output, "value") {
		t.Errorf("PrintProperty output should contain label and value, got: %q", output)
	}
}

func TestPrintPropertyToStderr(t *testing.T) {
	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	PrintPropertyToStderr("LABEL", "value")

	w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()
	os.Stderr = oldStderr

	if !strings.Contains(output, "LABEL") || !strings.Contains(output, "value") {
		t.Errorf("PrintPropertyToStderr output should contain label and value, got: %q", output)
	}
}

func TestPrintSuccess(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintSuccess("Success Title", "description")

	w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()
	os.Stdout = oldStdout

	if !strings.Contains(output, "STATUS") || !strings.Contains(output, "Success Title") || !strings.Contains(output, "description") {
		t.Errorf("PrintSuccess output should contain status, title, and description, got: %q", output)
	}
}

func TestPrintInfo(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintInfo("Info Title", "description")

	w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()
	os.Stdout = oldStdout

	if !strings.Contains(output, "INFO") || !strings.Contains(output, "Info Title") || !strings.Contains(output, "description") {
		t.Errorf("PrintInfo output should contain info, title, and description, got: %q", output)
	}
}

func TestPrintProperty_EmptyLabel(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintProperty("", "value")

	w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()
	os.Stdout = oldStdout

	// Should not print anything
	if output != "" {
		t.Errorf("PrintProperty with empty label should not print anything, got: %q", output)
	}
}
