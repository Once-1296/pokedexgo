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
			input : "   Charmander   Bulbasaur   PIKACHU   ",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input: " aWwAB TAha SaRaNg nathAN",
			expected: []string{"awwab","taha","sarang","nathan"},
		},
		{
			input: "   BooTDEV Is  so  CrAzY   ",
			expected: []string{"bootdev","is","so","crazy"},
		},
		// add more cases here
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice
		// if they don't match, use t.Errorf and continue to the next case
		if len(actual) != len(c.expected) {
			// error and continue here
			t.Errorf("Expected %d words but got %d words", len(c.expected), len(actual))
			t.Fail()
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord {
				t.Errorf("Expected %s at index %d, but got %s.", expectedWord, i,  word)
				t.Fail()
			}
		}
	}
}
