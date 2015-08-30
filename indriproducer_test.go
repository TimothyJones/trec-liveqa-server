package main

import (
	"testing"
)

func TestPrepareSDQuery(t *testing.T) {
	original := "What is artificial intelligence?"
	expected := "#weight( 0.85 #combine( what is artificial intelligence ) 0.10 #combine( #1( what is ) #1( is artificial ) #1( artificial intelligence ) ) 0.05 #combine ( #uw8( what is ) #uw8( is artificial ) #uw8( artificial intelligence ) ) )"

	if sd := PrepareSDQuery(GetQueryTerms(original)); sd != expected {
		t.Errorf("Expected payload `%s` but got `%s`", expected, sd)
	}
}

func TestRemoveDuplicateDocnos(t *testing.T) {
	result := RemoveDuplicateDocnos([]string{"1", "2", "3", "2", "3", "4", "5", "2", "1", "4"})
	expected := []string{"1", "2", "3", "4", "5"}

	if len(expected) != len(result) {
		t.Errorf("Expected `%s` but got `%s`", expected, result)
	} else {
		mismatch := false
		for i := range result {
			if expected[i] != result[i] {
				mismatch = true
			}
		}
		if mismatch {
			t.Errorf("Expected `%s` but got `%s`", expected, result)
		}
	}
}
