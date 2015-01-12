package gec

type Extractor struct {
	tp *TextProcessor
	bp *BlockProcessor
}

func NewExtractor(o *Option) (e *Extractor) {
	tp := NewTextProcessor(o)
	bp := NewBlockProcessor(o, tp)
	return &Extractor{tp: tp, bp: bp}
}

func (self *Extractor) ExtractTitle(doc string) (t string) {
	if self.tp.HasFramesetOrRedirect(doc) {
		t = self.tp.ParseTitle(doc)
	} else {
		head := self.tp.ParseHeadHTML(doc)
		t = self.tp.ParseTitle(head)
	}
	return
}

func (self *Extractor) ExtractContent(doc string) (c string) {
	if self.tp.HasFramesetOrRedirect(doc) {
		c = ""
	} else {
		body := self.tp.ParseBodyHTML(doc)
		body = self.tp.ParseGoogleAdsSectionTargetHTML(body)
		body = self.tp.EliminateUselessTags(body)

		title := self.tp.ParseTitle(doc)
		body = self.tp.ReplaceHTag(body, title)

		for _, b := range self.tp.ParseBlock(body) {
			self.bp.Process(b)
		}
		c = self.bp.GetMaxScoreContent()
	}
	return
}
