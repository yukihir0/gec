package gec

func Analyse(doc string, o *Option) (content, title string) {
	opt := NewOption()
	if o != nil {
		opt = o
	}

	e := NewExtractor(opt)
	title = e.ExtractTitle(doc)
	content = e.ExtractContent(doc)
	return
}
