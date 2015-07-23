package gec

import (
	"regexp"
)

// Option represents extract options
type Option struct {
	Threashold        float64        // 本文と見なすスコアの閾値
	MinLength         int            // 評価を行うブロック長の最小値
	DecayFactor       float64        // 減衰係数(小さいほど先頭に近いブロックのスコアが高くなる)
	ContinuousFactor  float64        // 連続ブロック係数(大きいほどブロックを連続と判定しにくくなる)
	NotBodyFactor     float64        // 非Body係数(大きいほどブロックのスコアが高くなる)
	PunctuationWeight int            // 句読点に対するスコア
	Punctuations      *regexp.Regexp // 句読点
	WasteExpressions  *regexp.Regexp // フッタに含まれる特徴的なキーワード
	DomSeparator      string         // DOM間に挿入する文字列
	Debug             bool           // ブロック情報を出力
}

func (o *Option) initialize() {
	o.Threashold = 100.0
	o.MinLength = 80
	o.DecayFactor = 0.73
	o.ContinuousFactor = 1.62
	o.NotBodyFactor = 0.72
	o.PunctuationWeight = 10
	o.Punctuations = regexp.MustCompile("([、。，．！？]|\\.[^A-Za-z0-9]|,[^0-9]|!|\\?)")
	o.WasteExpressions = regexp.MustCompile("(?i)Copyright|All Rights Reserved")
	o.DomSeparator = ""
	o.Debug = false
}

// NewOption initialize option
func NewOption() (o *Option) {
	o = &Option{}
	o.initialize()
	return
}
