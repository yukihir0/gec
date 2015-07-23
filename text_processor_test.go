package gec

import (
	"fmt"
	"testing"
)

func TestZenkakuASCII2HankakuASCIILowerAlphabet(t *testing.T) {
	opt := NewOption()
	tp := newTextProcessor(opt)

	actual := tp.zenkaku2Hankaku("ａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ")
	expected := "abcdefghijklmnopqrstuvwxyz"

	if actual != expected {
		t.Errorf("expected %v, but got %v", expected, actual)
	}
}

func TestZenkakuASCII2HankakuASCIIUpperAlphabet(t *testing.T) {
	opt := NewOption()
	tp := newTextProcessor(opt)

	actual := tp.zenkaku2Hankaku("ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺ")
	expected := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	if actual != expected {
		t.Errorf("expected %v, but got %v", expected, actual)
	}
}

func TestZenkakuASCII2HankakuASCIINumber(t *testing.T) {
	opt := NewOption()
	tp := newTextProcessor(opt)

	actual := tp.zenkaku2Hankaku("０１２３４５６７８９")
	expected := "0123456789"

	if actual != expected {
		t.Errorf("expected %v, but got %v", expected, actual)
	}
}

func TestZenkakuASCII2HankakuASCIISymbol(t *testing.T) {
	opt := NewOption()
	tp := newTextProcessor(opt)

	actual := tp.zenkaku2Hankaku("：；＜＝＞？＠［＼］＾＿｀｛｜｝～！＂＃＄％＆＇（）＊，－．／a")
	expected := ":;<=>?@[\\]^_`{|}~!\"#$%&'()*,-./a"

	if actual != expected {
		t.Errorf("expected %v, but got %v", expected, actual)
	}
}

func TestZenkakuSpace2HankakuSpace(t *testing.T) {
	opt := NewOption()
	tp := newTextProcessor(opt)

	actual := tp.zenkaku2Hankaku("　")
	expected := " "

	if actual != expected {
		t.Errorf("expected %v, but got %v", expected, actual)
	}
}

func TestCharref2Ascii(t *testing.T) {
	opt := NewOption()
	tp := newTextProcessor(opt)

	actual := tp.charref2Ascii("&nbsp;&lt;&gt;&amp;&laquo;&raquo;")
	expected := " <>&«»"

	if actual != expected {
		t.Errorf("expected %v, but got %v", expected, actual)
	}
}

func TestRemoveKeisen(t *testing.T) {
	opt := NewOption()
	tp := newTextProcessor(opt)

	for r := 0x2500; r <= 0x257F; r++ {
		actual := tp.eliminateRuledLine(fmt.Sprintf("%c", r))
		expected := ""

		if actual != expected {
			t.Errorf("expected %v, but got %v", expected, actual)
		}
	}
}

func TestEliminateUselessTags(t *testing.T) {
	opt := NewOption()
	tp := newTextProcessor(opt)

	actuals := []string{
		tp.EliminateUselessTags("<script>hoge</script><style>fuga</style><select><option>a</option></select><noscript>foo</noscript>"),
		tp.EliminateUselessTags("<!-- comment\ncomment -->"),
		tp.EliminateUselessTags("<!abc 123>"),
		tp.EliminateUselessTags("<div class=\"alpslab-slide\">hoge</div>"),
		tp.EliminateUselessTags("<div id=\"more\">hoge</div>"),
		tp.EliminateUselessTags("<div class=\"more\">hoge</div>"),
	}
	expecteds := []string{
		"",
		"",
		"",
		"",
		"",
		"",
	}

	for i := range actuals {
		if actuals[i] != expecteds[i] {
			t.Errorf("expected %v, but got %v", expecteds[i], actuals[i])
		}
	}
}

func TestHasOnlyTags(t *testing.T) {
	opt := NewOption()
	tp := newTextProcessor(opt)

	actuals := []bool{
		tp.IsOnlyTags("\t <br> \n"),
		tp.IsOnlyTags("\t hoge \n"),
	}
	expecteds := []bool{
		true,
		false,
	}

	for i := range actuals {
		if actuals[i] != expecteds[i] {
			t.Errorf("expected %v, but got %v", expecteds[i], actuals[i])
		}
	}
}

func TestEliminateLink(t *testing.T) {
	opt := NewOption()
	tp := newTextProcessor(opt)

	actuals := []string{
		tp.EliminateLink("<div>content</div><li>1.hoge</li><li>2.fuga</li></ul>"),
		tp.EliminateLink("<div>content</div><a href=\"/hoge\">hoge</a><a href=\"/fuga\">fuga</a><ul><li>1.hoge</li><li>2.fuga</li></ul>"),
	}
	expecteds := []string{
		"content1.hoge2.fuga",
		"",
	}

	for i := range actuals {
		if actuals[i] != expecteds[i] {
			t.Errorf("expected %v, but got %v", expecteds[i], actuals[i])
		}
	}
}

func TestIsLinklist(t *testing.T) {
	opt := NewOption()
	tp := newTextProcessor(opt)

	actuals := []bool{
		tp.isLinkList("<ul><li>1.hoge</li><li>2.fuga</li></ul>"),
		tp.isLinkList("start 123<ul><li>1.hoge</li><li>2.fuga</li></ul>456 end"),
	}
	expecteds := []bool{
		true,
		false,
	}

	for i := range actuals {
		if actuals[i] != expecteds[i] {
			t.Errorf("expected %v, but got %v", expecteds[i], actuals[i])
		}
	}
}

func TestStripTags(t *testing.T) {
	opt := NewOption()
	tp := newTextProcessor(opt)

	actual := tp.EliminateTags("Ａ　ＢＣ\342\224\200<br>\n<br>&lt;&gt; \t\n ", "")
	expected := "A BC\n<> \n"

	if actual != expected {
		t.Errorf("expected %v, but got %v", expected, actual)
	}
}

/*
func TestEvaluatelist(t *testing.T) {
	actual := tp.caluculateLinkRate([]string{"hoge", "<a href=\"/hoge\">hoge</a>", "<a href=\"/fuga\">fuga</a>", "<a href=\"/foo\">foo</a>"})
	expected := 6.0625

	if actual != expected {
		t.Errorf("expected %v, but got %v", expected, actual)
	}
}
*/
