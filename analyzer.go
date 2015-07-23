package gec

// Analyse to extract content and title
func Analyse(doc string, o *Option) (content, title string) {
	opt := NewOption()
	if o != nil {
		opt = o
	}

	e := newExtractor(opt)
	title = e.ExtractTitle(doc)
	content = e.ExtractContent(doc)
	return
}
