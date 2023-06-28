package easyjson

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	jsonStr := `{"name":"wang","age":10,"like":{"sports":"football"}}`
	like, err := Get(jsonStr, `like/sports`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(like)
}
