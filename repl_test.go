package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "     CHARIZARD PIKACHU BULBASAUR squirtle     ",
			expected: []string{"charizard", "pikachu", "bulbasaur", "squirtle"},
		},
		// more cases go from here...
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("the test case has failed, because `len(actual) != len(c.expected)`\n\t%v != %v", len(actual), len(c.expected))
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("the test case has failed, because `word != expectedWord`\n\t%v != %v", word, expectedWord)
				break
			}
		}
	}
}
