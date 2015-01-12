package gec

import (
	"log"
	"math"
	"sort"
)

const (
	DECAY      = 1.0
	CONTINUOUS = 1.0
)

type BlockProcessor struct {
	o  *Option
	tp *TextProcessor
	bl BlockList
	df float64
	cf float64
}

func NewBlockProcessor(o *Option, tp *TextProcessor) (bp *BlockProcessor) {
	return &BlockProcessor{o: o, tp: tp, bl: NewBlockList(), df: DECAY, cf: CONTINUOUS}
}

func (self *BlockProcessor) Process(doc string) {
	if !self.isTargetBlock(doc) {
		return
	}

	bs, cs := self.calculateScore(doc)
	if self.o.Debug {
		log.Printf("[debug]block_score: %.3f, continuous_score: %.3f\n", bs, cs)
		log.Printf("[debug]block_length: %d\n", len(doc))
		log.Printf("[debug]%s\n", self.tp.EliminateTags(doc, ""))
	}

	self.updateDecay()
	if self.hasContinuousBlock() {
		self.updateContinuous()
	}

	if self.clusterBlock(doc, bs, cs) {
		self.resetContinuous()
	}
}

func (self *BlockProcessor) GetMaxScoreContent() (c string) {
	sort.Sort(self.bl)
	return self.tp.EliminateTags(self.bl[0].Text, self.o.DomSeparator)
}

func (self *BlockProcessor) isTargetBlock(doc string) (b bool) {
	return !(self.tp.IsZeroLength(doc) || self.tp.IsShortLength(doc) || self.tp.IsOnlyTags(doc))
}

func (self *BlockProcessor) calculateScore(doc string) (bs, cs float64) {
	notLinkedLength := len(self.tp.EliminateLink(doc))
	punctuationLength := len(self.tp.ParsePunctuations(doc))
	bs = float64(notLinkedLength+punctuationLength*self.o.PunctuationWeight) * self.df

	wasteLength := len(self.tp.ParseWasteExpressions(doc))
	amazonLength := len(self.tp.ParseAmazons(doc))
	notBodyRate := float64(wasteLength) + float64(amazonLength)/2.0
	if notBodyRate > 0.0 {
		bs *= math.Pow(self.o.NotBodyFactor, notBodyRate)
	}
	cs = bs * self.cf
	return
}

func (self *BlockProcessor) updateDecay() {
	self.df *= self.o.DecayFactor
}

func (self *BlockProcessor) updateContinuous() {
	self.cf /= self.o.ContinuousFactor
}

func (self *BlockProcessor) resetContinuous() {
	self.cf = CONTINUOUS
}

func (self *BlockProcessor) hasContinuousBlock() (b bool) {
	return len(self.bl[len(self.bl)-1].Text) > 0
}

func (self *BlockProcessor) clusterBlock(doc string, bs, cs float64) (b bool) {
	if cs > self.o.Threashold {
		self.bl[len(self.bl)-1].Text += doc + "\n"
		self.bl[len(self.bl)-1].Score += cs
		b = true
	} else if bs > self.o.Threashold {
		self.bl.AppendBlock()
		self.bl[len(self.bl)-1].Text = doc + "\n"
		self.bl[len(self.bl)-1].Score = bs
		b = true
	} else {
		b = false
	}
	return
}
