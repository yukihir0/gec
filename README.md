# gec

"gec" is port of ExtractContent.rb by golang.

## Original
- http://labs.cybozu.co.jp/blog/nakatani/2007/09/web_1.html

## Install

```
go get github.com/yukihir0/gec
```

## How to use

```
text := "..."
opt := gec.NewOption()
content, title := gec.Analyse(text, opt)
```

## License

Copyright &copy; 2015 yukihir0
