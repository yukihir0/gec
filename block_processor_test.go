package gec

import (
	"testing"
)

func TestMaxScoreContent(t *testing.T) {
	opt := NewOption()
	tp := NewTextProcessor(opt)
	bp := NewBlockProcessor(opt, tp)

	docs := []string{
		"",
		"<html><head><title>hoge</title></head><html>",
	}
	expecteds := []string{
		"",
		"",
	}

	for i, _ := range docs {
		bp.Process(docs[i])
		actual := bp.GetMaxScoreContent()
		if actual != expecteds[i] {
			t.Errorf("expected %s, but got %s", expecteds[i], actual)
		}
	}

}
