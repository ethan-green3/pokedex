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
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "         Wow----     This_*& is     A  Wacky String_ &          ",
			expected: []string{"wow----", "this_*&", "is", "a", "wacky", "string_", "&"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "                                   ",
			expected: []string{},
		},
		{
			input:    "\nCharmander\tPIKACHU",
			expected: []string{"charmander", "pikachu"},
		},
		{
			input:    "\t\n",
			expected: []string{},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(c.expected) != len(actual) {
			t.Errorf("Expected length %v does not equal actual length %v", len(c.expected), len(actual))
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Word %s does not equal expected word %s", word, expectedWord)
			}
		}
	}
}
