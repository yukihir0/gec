package gec

import (
	"html"
	"math"
	"regexp"
	"strings"
	"unicode/utf8"
)

type textProcessor struct {
	o  *Option
	rm map[string]*regexp.Regexp
}

func newTextProcessor(o *Option) (tp *textProcessor) {
	return &textProcessor{o: o, rm: initializeRegexpMap()}
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
		"heading":  regexp.MustCompile("(?i)(<h\\d\\s*>\\s*(.*?)\\s*<\\/h\\d\\s*>)"),
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

func (tp *textProcessor) HasFramesetOrRedirect(doc string) (b bool) {
	return tp.rm["frameset"].MatchString(doc)
}

func (tp *textProcessor) EliminateTags(doc string, separator string) (s string) {
	s = tp.rm["tags"].ReplaceAllString(doc, separator)
	s = tp.zenkaku2Hankaku(s)
	s = tp.eliminateRuledLine(s)
	s = tp.charref2Ascii(s)
	s = html.UnescapeString(s)
	s = tp.rm["tab"].ReplaceAllString(s, " ")
	s = tp.rm["cr"].ReplaceAllString(s, "\n")
	return
}

func (tp *textProcessor) EliminateUselessTags(doc string) (s string) {
	s = tp.rm["useless1"].ReplaceAllString(doc, "")
	s = tp.rm["useless2"].ReplaceAllString(s, "")
	s = tp.rm["useless3"].ReplaceAllString(s, "")
	s = tp.rm["useless4"].ReplaceAllString(s, "")
	s = tp.rm["useless5"].ReplaceAllString(s, "")
	s = tp.rm["useless6"].ReplaceAllString(s, "")
	return
}

func (tp *textProcessor) EliminateLink(doc string) (s string) {
	linkLength := len(tp.rm["a"].FindAllString(doc, -1))
	s = tp.rm["a"].ReplaceAllString(doc, "")
	s = tp.rm["form"].ReplaceAllString(s, "")
	s = tp.EliminateTags(s, "")
	if len(s) < 20*linkLength || tp.isLinkList(doc) {
		s = ""
	}
	return
}

func (tp *textProcessor) ParseHeadHTML(doc string) (s string) {
	head, _ := tp.splitHeadBody(doc)
	return head
}

func (tp *textProcessor) ParseBodyHTML(doc string) (s string) {
	_, body := tp.splitHeadBody(doc)
	return body
}

func (tp *textProcessor) ParseTitle(doc string) (s string) {
	t := tp.rm["title"].FindStringSubmatch(doc)
	if t != nil {
		s = tp.EliminateTags(t[1], "")
	} else {
		s = ""
	}
	return
}

func (tp *textProcessor) ParseGoogleAdsSectionTargetHTML(doc string) (s string) {
	s = tp.rm["gads1"].ReplaceAllString(doc, "")
	if tp.rm["gads2"].MatchString(s) {
		target := tp.rm["gads3"].FindAllString(s, -1)
		if target != nil {
			s = ""
			for _, b := range target {
				s += b + "\n"
			}
		}
	}
	return
}

func (tp *textProcessor) ReplaceHeadingTag(doc, title string) (s string) {
	s = doc
	for _, b := range tp.rm["heading"].FindAllStringSubmatch(s, -1) {
		if len(b[2]) >= 3 && strings.Contains(title, b[2]) {
			s = strings.Replace(s, b[1], "<div>"+b[2]+"</div>", -1)
		}
	}
	return
}

func (tp *textProcessor) ParsePunctuations(doc string) (s []string) {
	return tp.o.Punctuations.FindAllString(doc, -1)
}

func (tp *textProcessor) ParseWasteExpressions(doc string) (s []string) {
	return tp.o.WasteExpressions.FindAllString(doc, -1)
}

func (tp *textProcessor) ParseAmazons(doc string) (s []string) {
	return tp.rm["amazon"].FindAllString(doc, -1)
}

func (tp *textProcessor) ParseBlock(doc string) (s []string) {
	for _, l := range tp.rm["body"].Split(doc, -1) {
		s = append(s, strings.TrimSpace(l))
	}
	return
}

func (tp *textProcessor) IsZeroLength(doc string) (b bool) {
	return len(doc) == 0
}

func (tp *textProcessor) IsShortLength(doc string) (b bool) {
	return len(tp.EliminateLink(doc)) < tp.o.MinLength
}

func (tp *textProcessor) eliminateTag(doc string) (s string) {
	s = tp.rm["tag"].ReplaceAllString(doc, "")
	s = tp.rm["nbsp"].ReplaceAllString(s, "")
	s = strings.TrimSpace(s)
	return s
}

func (tp *textProcessor) IsOnlyTags(doc string) (b bool) {
	return len(tp.eliminateTag(doc)) == 0
}

func (tp *textProcessor) splitHeadBody(doc string) (head, body string) {
	loc := tp.rm["head"].FindStringIndex(doc)
	if loc != nil {
		head = doc[0:loc[0]]
		body = doc[loc[1]:]
	} else {
		head = doc
		body = doc
	}
	return
}

func (tp *textProcessor) isLinkList(doc string) (b bool) {
	ul := tp.rm["uldlol"].FindStringSubmatch(doc)
	if ul != nil {
		li := tp.rm["li"].Split(ul[1], -1)
		li = li[1:]
		outside := tp.rm["uldl"].ReplaceAllString(doc, "")
		outside = tp.rm["tags"].ReplaceAllString(outside, "")
		outside = tp.rm["space"].ReplaceAllString(outside, " ")

		rate := tp.caluculateLinkRate(li)
		b = float64(len(outside)) <= (float64(len(doc)) * rate)
	} else {
		b = false
	}
	return
}

func (tp *textProcessor) caluculateLinkRate(doc []string) (r float64) {
	if len(doc) == 0 {
		return 1.0
	}

	hit := 0.0
	for _, li := range doc {
		if tp.rm["href"].MatchString(li) {
			hit += 1.0
		}
	}
	return (9.0*math.Pow(hit/float64(len(doc)), 2) + 1.0) / 45.0
}

func (tp *textProcessor) zenkaku2Hankaku(doc string) (s string) {
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

func (tp *textProcessor) eliminateRuledLine(doc string) (s string) {
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

func (tp *textProcessor) charref2Ascii(doc string) (s string) {
	r := strings.NewReplacer("&nbsp;", " ", "&lt;", "<", "&gt;", ">", "&amp;", "&", "&laquo;", "\xc2\xab", "&raquo;", "\xc2\xbb")
	return r.Replace(doc)
}
