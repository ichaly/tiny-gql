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
	l := NewLexer(`"""abc\n
	   123
	\\1111\"\"🤔\t\n\r\f
\u3d7d
    22222"""`)
	for tok, err := l.ReadToken(); tok.Kind != EOF; tok, err = l.ReadToken() {
		if err == nil {
			println(tok.Kind, tok.Value)
		} else {
			println(err.Error())
			break
		}
	}
}

func TestReadString(t *testing.T) {
	l := NewLexer(`"\u3d7d 这事一个字符串的测试🤔\n\t\r"`)
	for tok, err := l.ReadToken(); tok.Kind != EOF; tok, err = l.ReadToken() {
		if err != nil {
			println(err.Error())
			break
		} else {
			println(tok.Kind, tok.Value)
		}
	}
}
