package main

import "testing"

func TestUrl(t *testing.T) {
	query := "query words with a question?"
	var TestUrlTests = []struct {
		host     string
		expected string
	}{
		{"http://localhost:8080/", "http://localhost:8080/searchxml?n=20&q=query+words+with+a+question%3F"},
		{"http://localhost:8080", "http://localhost:8080/searchxml?n=20&q=query+words+with+a+question%3F"},
		{"http://localhost:8080/foo", "http://localhost:8080/foo/searchxml?n=20&q=query+words+with+a+question%3F"},
		{"http://a.b.c.com/foo", "http://a.b.c.com/foo/searchxml?n=20&q=query+words+with+a+question%3F"},
		{"https://a.b.c.com/foo", "https://a.b.c.com/foo/searchxml?n=20&q=query+words+with+a+question%3F"},
		{"a.b.c.com", "http://a.b.c.com/searchxml?n=20&q=query+words+with+a+question%3F"},
	}

	for _, test := range TestUrlTests {
		ap := &GalagoAnswerProducer{Host: test.host}

		url, err := ap.GetUrl(query)
		if err != nil {
			t.Errorf("Unexpected error '%s'", err.Error())
			continue
		}
		actual := url.String()

		if test.expected != actual {
			t.Errorf("Expected '%s' but got '%s'", test.expected, actual)
		}
	}

}
