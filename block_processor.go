package gec

import (
	"log"
	"math"
	"sort"
)

// threshold for extract block
const (
	Decay      = 1.0
	Continuous = 1.0
)

type blockProcessor struct {
	o  *Option
	tp *textProcessor
	bl blockList
	df float64
	cf float64
}

func newBlockProcessor(o *Option, tp *textProcessor) (bp *blockProcessor) {
	return &blockProcessor{o: o, tp: tp, bl: newBlockList(), df: Decay, cf: Continuous}
}

func (bp *blockProcessor) Cluster(doc string) {
	if !bp.isCandidate(doc) {
		return
	}

	bs, cs := bp.calculateScore(doc)
	if bp.o.Debug {
		log.Printf("[debug]block_score: %.3f, continuous_score: %.3f\n", bs, cs)
		log.Printf("[debug]block_length: %d\n", len(doc))
		log.Printf("[debug]%s\n", bp.tp.EliminateTags(doc, ""))
	}

	bp.updateDecay()
	if bp.hasContinuousBlock() {
		bp.updateContinuous()
	}

	if bp.clusterBlock(doc, bs, cs) {
		bp.resetContinuous()
	}
}

func (bp *blockProcessor) GetMaxScoreContent() (c string) {
	sort.Sort(bp.bl)
	return bp.tp.EliminateTags(bp.bl[0].Text, bp.o.DomSeparator)
}

func (bp *blockProcessor) isCandidate(doc string) (b bool) {
	return !(bp.tp.IsZeroLength(doc) || bp.tp.IsShortLength(doc) || bp.tp.IsOnlyTags(doc))
}

func (bp *blockProcessor) calculateScore(doc string) (bs, cs float64) {
	notLinkedLength := len(bp.tp.EliminateLink(doc))
	punctuationLength := len(bp.tp.ParsePunctuations(doc))
	bs = float64(notLinkedLength+punctuationLength*bp.o.PunctuationWeight) * bp.df

	wasteLength := len(bp.tp.ParseWasteExpressions(doc))
	amazonLength := len(bp.tp.ParseAmazons(doc))
	notBodyRate := float64(wasteLength) + float64(amazonLength)/2.0
	if notBodyRate > 0.0 {
		bs *= math.Pow(bp.o.NotBodyFactor, notBodyRate)
	}
	cs = bs * bp.cf
	return
}

func (bp *blockProcessor) updateDecay() {
	bp.df *= bp.o.DecayFactor
}

func (bp *blockProcessor) updateContinuous() {
	bp.cf /= bp.o.ContinuousFactor
}

func (bp *blockProcessor) resetContinuous() {
	bp.cf = Continuous
}

func (bp *blockProcessor) hasContinuousBlock() (b bool) {
	return len(bp.bl[len(bp.bl)-1].Text) > 0
}

func (bp *blockProcessor) clusterBlock(doc string, bs, cs float64) (b bool) {
	if cs > bp.o.Threashold {
		bp.bl[len(bp.bl)-1].Text += doc + "\n"
		bp.bl[len(bp.bl)-1].Score += cs
		b = true
	} else if bs > bp.o.Threashold {
		bp.bl.AppendBlock()
		bp.bl[len(bp.bl)-1].Text = doc + "\n"
		bp.bl[len(bp.bl)-1].Score = bs
		b = true
	} else {
		b = false
	}
	return
}
