package lexer

import (
	"github.com/bytedance/sonic"
	"testing"
)

func TestReadToken(t *testing.T) {
	gql := `
#	查询粉丝的访问痕迹

query ($orderStr: String, $targetType: Int) {
  getFollowUsersTraces(orderStr: $orderStr, targetType: $targetType) {
	orderStr
	unReadFollowUserTraceNum
  }
}
`
	var tokens []Token
	l := New(&Input{Content: gql})
	for {
		tok, err := l.ReadToken()
		if err != nil {
			break
		}
		if tok.Kind == EOF {
			break
		}
		tokens = append(tokens, tok)
	}
	output, err := sonic.Marshal(&tokens)
	if err != nil {
		return
	}
	println(string(output))
}
