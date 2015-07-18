package gec

import (
	"testing"
)

func TestAnalyse(t *testing.T) {
	docs := []string{
		"<html><head><title>hoge</title></head><frameset></frameset></html>",
		"<html><head><title>hoge</title><meta http-equiv=\"refresh\" content=\"5;URL=http://www.example.com\"></head></html>",
		"<html><head><title>fuga</title></head><body></body></html>",
		"<html><head><title>fuga</title></head><body><!-- google_ad_section_start -->a<!-- google_ad_section_end --><!-- google_ad_section_start -->b<!-- google_ad_section_end --></body></html>",
		"<html><head><title>fuga</title></head><body><h1>fuga</h1><h2>fuga</h2></body></html>",
	}
	expecteds := [][]string{
		[]string{"", "hoge"},
		[]string{"", "hoge"},
		[]string{"", "fuga"},
		[]string{"", "fuga"},
		[]string{"", "fuga"},
	}

	for i, _ := range docs {
		content, title := Analyse(docs[i], nil)
		if content != expecteds[i][0] {
			t.Errorf("expected %v, but got %v", expecteds[i][0], content)
		}

		if title != expecteds[i][1] {
			t.Errorf("expected %v, but got %v", expecteds[i][1], title)
		}
	}
}

func TestAnalyseWithOption(t *testing.T) {
	opt := NewOption()
	docs := []string{
		"<html><head><title>hoge</title></head><frameset></frameset></html>",
		"<html><head><title>hoge</title><meta http-equiv=\"refresh\" content=\"5;URL=http://www.example.com\"></head></html>",
		"<html><head><title>fuga</title></head><body></body></html>",
		"<html><head><title>fuga</title></head><body><!-- google_ad_section_start -->a<!-- google_ad_section_end --><!-- google_ad_section_start -->b<!-- google_ad_section_end --></body></html>",
		"<html><head><title>fuga</title></head><body><h1>fuga</h1><h2>fuga</h2></body></html>",
	}
	expecteds := [][]string{
		[]string{"", "hoge"},
		[]string{"", "hoge"},
		[]string{"", "fuga"},
		[]string{"", "fuga"},
		[]string{"", "fuga"},
	}

	for i, _ := range docs {
		content, title := Analyse(docs[i], opt)
		if content != expecteds[i][0] {
			t.Errorf("expected %v, but got %v", expecteds[i][0], content)
		}

		if title != expecteds[i][1] {
			t.Errorf("expected %v, but got %v", expecteds[i][1], title)
		}
	}
}
