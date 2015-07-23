package gec

type extractor struct {
	tp *textProcessor
	bp *blockProcessor
}

func newExtractor(o *Option) (e *extractor) {
	tp := newTextProcessor(o)
	bp := newBlockProcessor(o, tp)
	return &extractor{tp: tp, bp: bp}
}

func (e *extractor) ExtractTitle(doc string) (t string) {
	if e.tp.HasFramesetOrRedirect(doc) {
		t = e.tp.ParseTitle(doc)
	} else {
		head := e.tp.ParseHeadHTML(doc)
		t = e.tp.ParseTitle(head)
	}
	return
}

func (e *extractor) ExtractContent(doc string) (c string) {
	if e.tp.HasFramesetOrRedirect(doc) {
		c = ""
	} else {
		body := e.tp.ParseBodyHTML(doc)
		body = e.tp.ParseGoogleAdsSectionTargetHTML(body)
		body = e.tp.EliminateUselessTags(body)

		title := e.tp.ParseTitle(doc)
		body = e.tp.ReplaceHeadingTag(body, title)

		for _, b := range e.tp.ParseBlock(body) {
			e.bp.Cluster(b)
		}
		c = e.bp.GetMaxScoreContent()
	}
	return
}
