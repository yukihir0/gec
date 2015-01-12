package gec

type Block struct {
	Text  string
	Score float64
}

func NewBlock(text string, score float64) (b *Block) {
	return &Block{Text: text, Score: score}
}
