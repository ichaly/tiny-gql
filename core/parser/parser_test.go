package parser

import (
	"github.com/bytedance/sonic"
	"github.com/ichaly/tiny-go/core/lexer"
	"testing"
)

func TestParseQuery(t *testing.T) {
	gql := `
#	查询粉丝的访问痕迹

query ($orderStr: String, $targetType: Int) {
  getFollowUsersTraces(orderStr: $orderStr, targetType: $targetType) {
	orderStr
	unReadFollowUserTraceNum
  }
}
`
	query, err := ParseQuery(&lexer.Input{Content: gql})
	if err != nil {
		return
	}
	output, err := sonic.Marshal(&query)
	if err != nil {
		return
	}
	println(string(output))
}
