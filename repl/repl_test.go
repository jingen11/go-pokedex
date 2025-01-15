package repl

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		}, {
			input:    "hello, world!",
			expected: []string{"hello,", "world!"},
		}, {
			input:    "This IS a CapiTal LeTTER",
			expected: []string{"this", "is", "a", "capital", "letter"},
		},
	}

	for _, c := range cases {
		actual := CleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("Expected length: %d, Get: %d", len(c.expected), len(actual))
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("Expected: %s, Get: %s", expectedWord, word)
			}
		}
	}
}
