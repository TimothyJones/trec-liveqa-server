package main

import (
	"testing"
)

func TestPrepareSDQuery(t *testing.T) {
	original := "What is artificial intelligence?"
	expected := "#weight( 0.85 #combine( what is artificial intelligence ) 0.10 #combine( #1( artificial intelligence ) ) 0.05 #combine ( #uw8( artificial intelligence ) ) )"

	if sd := PrepareSDQuery(GetQueryTerms(original)); sd != expected {
		t.Errorf("Expected payload `%s` but got `%s`", expected, sd)
	}
}
