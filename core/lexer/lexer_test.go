package lexer

import (
	"testing"
)

const gql = `




 	#	查询粉丝的访问痕迹

query ($orderStr: String, $targetType: Int) {
  getFollowUsersTraces(orderStr: "desc", targetType: $targetType) {
	orderStr
	unReadFollowUserTraceNum
  }
}
`

func TestTrim(t *testing.T) {
	l := NewLexer(gql)
	l.trim()
	println(l.text[l.start:l.end], l.start, l.end, l.line)
}

func TestReadBlock(t *testing.T) {
	l := NewLexer(`"""abc\n""
	   123"
	\\1111
	\"""
\u3d7d
    22222"""`)
	for tok, err := l.ReadToken(); tok.Kind != EOF && err == nil; tok, err = l.ReadToken() {
		println(tok.Kind, tok.Value)
	}
}

func TestReadString(t *testing.T) {
	l := NewLexer(`"\u3d7d"`)
	for tok, err := l.ReadToken(); tok.Kind != EOF; tok, err = l.ReadToken() {
		if err != nil {
			println(err.Error())
		} else {
			println(tok.Kind, tok.Value)
		}
	}
}
