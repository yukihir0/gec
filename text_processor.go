package gec

import (
	"html"
	"math"
	"regexp"
	"strings"
	"unicode/utf8"
)

type TextProcessor struct {
	o  *Option
	rm map[string]*regexp.Regexp
}

func NewTextProcessor(o *Option) (tp *TextProcessor) {
	return &TextProcessor{o: o, rm: initializeRegexpMap()}
}

func initializeRegexpMap() (rm map[string]*regexp.Regexp) {
	rm = map[string]*regexp.Regexp{
		"frameset": regexp.MustCompile("(?i)<\\/frameset>|<meta\\s+http-equiv\\s*=\\s*[\"']?refresh['\"]?[^>]*url"),
		"head":     regexp.MustCompile("(?is)<\\/head\\s*>"),
		"title":    regexp.MustCompile("(?i)<title[^>]*>\\s*(.*?)\\s*<\\/title\\s*>"),
		"body":     regexp.MustCompile("<\\/?(?:div|center|td)[^>]*>|<p\\s*[^>]*class\\s*=\\s*[\"']?(?:posted|plugin-\\w+)['\"]?[^>]*>"),
		"nbsp":     regexp.MustCompile("&nbsb;"),
		"tab":      regexp.MustCompile("[ \t]+"),
		"space":    regexp.MustCompile("\\s+"),
		"cr":       regexp.MustCompile("\n\\s*"),
		"tag":      regexp.MustCompile("(?is)<[^>]*>"),
		"tags":     regexp.MustCompile("(?s)<.+?>"),
		"h":        regexp.MustCompile("(?i)(<h\\d\\s*>\\s*(.*?)\\s*<\\/h\\d\\s*>)"),
		"a":        regexp.MustCompile("(?is)<a\\s[^>]*>.*?<\\/a\\s*>"),
		"href":     regexp.MustCompile("(?is)<a\\s+href=(['\"]?)([^\"'\\s]+)"),
		"form":     regexp.MustCompile("(?is)<form\\s[^>]*>.*?<\\/form\\s*>"),
		"uldlol":   regexp.MustCompile("<(?:ul|dl|ol)(.+?)<\\/(?:ul|dl|ol)>"),
		"uldl":     regexp.MustCompile("(?is)<(?:ul|dl)(.+?)<\\/(?:ul|dl)>"),
		"li":       regexp.MustCompile("<li[^>]*>"),
		"gads1":    regexp.MustCompile("(?s)<!--\\s*google_ad_section_start\\(weight=ignore\\)\\s*-->.*?<!--\\s*google_ad_section_end.*?-->"),
		"gads2":    regexp.MustCompile("<!--\\s*google_ad_section_start[^>]*-->"),
		"gads3":    regexp.MustCompile("(?s)<!--\\s*google_ad_section_start[^>]*-->.*?<!--\\s*google_ad_section_end.*?-->"),
		"amazon":   regexp.MustCompile("(?i)amazon[a-z0-9\\.\\/\\-\\?&]+-22"),
		"useless1": regexp.MustCompile("[\342\200\230-\342\200\235]|[\342\206\220-\342\206\223]|[\342\226\240-\342\226\275]|[\342\227\206-\342\227\257]|\342\230\205|\342\230\206/"),
		"useless2": regexp.MustCompile("(?is)<(script|style|select|noscript)[^>]*>.*?<\\/(script|style|select|noscript)\\s*>"),
		"useless3": regexp.MustCompile("(?s)<!--.*?-->"),
		"useless4": regexp.MustCompile("<![A-Za-z].*?>"),
		"useless5": regexp.MustCompile("(?s)<div\\s[^>]*class\\s*=\\s*['\"]?alpslab-slide[\"']?[^>]*>.*?<\\/div\\s*>"),
		"useless6": regexp.MustCompile("(?i)<div\\s[^>]*(id|class)\\s*=\\s*['\"]?\\S*more\\S*[\"']?[^>]*>"),
	}
	return
}

func (self *TextProcessor) HasFramesetOrRedirect(doc string) (b bool) {
	return self.rm["frameset"].MatchString(doc)
}

func (self *TextProcessor) EliminateTag(doc string) (s string) {
	s = self.rm["tag"].ReplaceAllString(doc, "")
	s = self.rm["nbsp"].ReplaceAllString(s, "")
	s = strings.TrimSpace(s)
	return s
}

func (self *TextProcessor) EliminateTags(doc string, separator string) (s string) {
	s = self.rm["tags"].ReplaceAllString(doc, separator)
	s = self.zenkaku2Hankaku(s)
	s = self.eliminateRuledLine(s)
	s = self.charref2Ascii(s)
	s = html.UnescapeString(s)
	s = self.rm["tab"].ReplaceAllString(s, " ")
	s = self.rm["cr"].ReplaceAllString(s, "\n")
	return
}

func (self *TextProcessor) EliminateUselessTags(doc string) (s string) {
	s = self.rm["useless1"].ReplaceAllString(doc, "")
	s = self.rm["useless2"].ReplaceAllString(s, "")
	s = self.rm["useless3"].ReplaceAllString(s, "")
	s = self.rm["useless4"].ReplaceAllString(s, "")
	s = self.rm["useless5"].ReplaceAllString(s, "")
	s = self.rm["useless6"].ReplaceAllString(s, "")
	return
}

func (self *TextProcessor) EliminateLink(doc string) (s string) {
	linkLength := len(self.rm["a"].FindAllString(doc, -1))
	s = self.rm["a"].ReplaceAllString(doc, "")
	s = self.rm["form"].ReplaceAllString(s, "")
	s = self.EliminateTags(s, "")
	if len(s) < 20*linkLength || self.isLinkList(doc) {
		s = ""
	}
	return
}

func (self *TextProcessor) ParseHeadHTML(doc string) (s string) {
	head, _ := self.splitHeadBody(doc)
	return head
}

func (self *TextProcessor) ParseBodyHTML(doc string) (s string) {
	_, body := self.splitHeadBody(doc)
	return body
}

func (self *TextProcessor) ParseTitle(doc string) (s string) {
	t := self.rm["title"].FindStringSubmatch(doc)
	if t != nil {
		s = self.EliminateTags(t[1], "")
	} else {
		s = ""
	}
	return
}

func (self *TextProcessor) ParseGoogleAdsSectionTargetHTML(doc string) (s string) {
	s = self.rm["gads1"].ReplaceAllString(doc, "")
	if self.rm["gads2"].MatchString(s) {
		target := self.rm["gads3"].FindAllString(s, -1)
		if target != nil {
			s = ""
			for _, b := range target {
				s += b + "\n"
			}
		}
	}
	return
}

func (self *TextProcessor) ReplaceHTag(doc, title string) (s string) {
	s = doc
	for _, b := range self.rm["h"].FindAllStringSubmatch(s, -1) {
		if len(b[2]) >= 3 && strings.Contains(title, b[2]) {
			s = strings.Replace(s, b[1], "<div>"+b[2]+"</div>", -1)
		}
	}
	return
}

func (self *TextProcessor) ParsePunctuations(doc string) (s []string) {
	return self.o.Punctuations.FindAllString(doc, -1)
}

func (self *TextProcessor) ParseWasteExpressions(doc string) (s []string) {
	return self.o.WasteExpressions.FindAllString(doc, -1)
}

func (self *TextProcessor) ParseAmazons(doc string) (s []string) {
	return self.rm["amazon"].FindAllString(doc, -1)
}

func (self *TextProcessor) ParseBlock(doc string) (s []string) {
	for _, l := range self.rm["body"].Split(doc, -1) {
		s = append(s, strings.TrimSpace(l))
	}
	return
}

func (self *TextProcessor) IsZeroLength(doc string) (b bool) {
	return len(doc) == 0
}

func (self *TextProcessor) IsShortLength(doc string) (b bool) {
	return len(self.EliminateLink(doc)) < self.o.MinLength
}

func (self *TextProcessor) IsOnlyTags(doc string) (b bool) {
	return len(self.EliminateTag(doc)) == 0
}

func (self *TextProcessor) splitHeadBody(doc string) (head, body string) {
	loc := self.rm["head"].FindStringIndex(doc)
	if loc != nil {
		head = doc[0:loc[0]]
		body = doc[loc[1]:]
	} else {
		head = doc
		body = doc
	}
	return
}

func (self *TextProcessor) isLinkList(doc string) (b bool) {
	ul := self.rm["uldlol"].FindStringSubmatch(doc)
	if ul != nil {
		li := self.rm["li"].Split(ul[1], -1)
		li = li[1:]
		outside := self.rm["uldl"].ReplaceAllString(doc, "")
		outside = self.rm["tags"].ReplaceAllString(outside, "")
		outside = self.rm["space"].ReplaceAllString(outside, " ")

		rate := self.caluculateLinkRate(li)
		b = float64(len(outside)) <= (float64(len(doc)) * rate)
	} else {
		b = false
	}
	return
}

func (self *TextProcessor) caluculateLinkRate(doc []string) (r float64) {
	if len(doc) == 0 {
		return 1.0
	}

	hit := 0.0
	for _, li := range doc {
		if self.rm["href"].MatchString(li) {
			hit += 1.0
		}
	}
	return (9.0*math.Pow(hit/float64(len(doc)), 2) + 1.0) / 45.0
}

func (self *TextProcessor) zenkaku2Hankaku(doc string) (s string) {
	buf := make([]rune, 0, utf8.RuneCountInString(doc))
	for _, r := range doc {
		switch {
		// !"#$%&'()*+,-./
		case (r >= 0xFF01 && r <= 0xFF0F):
			fallthrough
		// 0-9
		case (r >= 0xFF10 && r <= 0xff19):
			fallthrough
		// :;<=>?@
		case (r >= 0xFF1A && r <= 0xff20):
			fallthrough
		// A-Z
		case (r >= 0xFF21 && r <= 0xFF3A):
			fallthrough
		// [\]^_`
		case (r >= 0xFF3B && r <= 0xFF40):
			fallthrough
		// a-z
		case (r >= 0xFF41 && r <= 0xFF5A):
			fallthrough
		// {|}~
		case (r >= 0xFF5B && r <= 0xFF5E):
			buf = append(buf, r-0xFEE0)
		// space
		case r == 0x3000:
			buf = append(buf, 0x0020)
		default:
			buf = append(buf, r)
		}
	}
	return string(buf)
}

func (self *TextProcessor) eliminateRuledLine(doc string) (s string) {
	buf := make([]rune, 0, utf8.RuneCountInString(doc))
	for _, r := range doc {
		switch {
		// ruled line
		case (string(r) >= "\342\224\200" && string(r) <= "\342\224\277"):
		case (string(r) >= "\342\225\200" && string(r) <= "\342\225\277"):
		default:
			buf = append(buf, r)
		}
	}
	return string(buf)
}

func (self *TextProcessor) charref2Ascii(doc string) (s string) {
	r := strings.NewReplacer("&nbsp;", " ", "&lt;", "<", "&gt;", ">", "&amp;", "&", "&laquo;", "\xc2\xab", "&raquo;", "\xc2\xbb")
	return r.Replace(doc)
}
