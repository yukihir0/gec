package gec

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestAnalyse(t *testing.T) {
	inputs := []string{
		"./data/input01.html",
		"./data/input02.html",
		"./data/input03.html",
		"./data/input04.html",
		"./data/input05.html",
		"./data/input06.html",
		"./data/input07.html",
		"./data/input08.html",
	}

	docs := []string{}
	for _, input := range inputs {
		html, _ := ioutil.ReadFile(input)
		docs = append(docs, string(html))
	}

	contents := []string{
		"./data/expected_content01.html",
		"./data/expected_content02.html",
		"./data/expected_content03.html",
		"./data/expected_content04.html",
		"./data/expected_content05.html",
		"./data/expected_content06.html",
		"./data/expected_content07.html",
		"./data/expected_content08.html",
	}
	titles := []string{
		"./data/expected_title01.html",
		"./data/expected_title02.html",
		"./data/expected_title03.html",
		"./data/expected_title04.html",
		"./data/expected_title05.html",
		"./data/expected_title06.html",
		"./data/expected_title07.html",
		"./data/expected_title08.html",
	}
	expecteds := [][]string{}
	for i := range contents {
		content, _ := ioutil.ReadFile(contents[i])
		title, _ := ioutil.ReadFile(titles[i])
		expecteds = append(expecteds, []string{string(content), strings.TrimRight(string(title), "\n")})
	}

	opt := NewOption()
	for i := range docs {
		content, title := Analyse(docs[i], opt)
		if content != expecteds[i][0] {
			t.Errorf("expected %s, but got %s", expecteds[i][0], content)
		}

		if title != expecteds[i][1] {
			t.Errorf("expected %s, but got %s", expecteds[i][1], title)
		}
	}
}
