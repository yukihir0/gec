package gec

type block struct {
	Text  string
	Score float64
}

func newBlock(text string, score float64) (b *block) {
	return &block{Text: text, Score: score}
}
