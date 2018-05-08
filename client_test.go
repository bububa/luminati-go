package luminati

import (
	"fmt"
	"testing"
)

func TestLuminatiClient_ShoulReturnContent_When200(t *testing.T) {
	c := NewClient("lum-customer-xxx-zone-residential", "xxx")

	c.SetFailuresLimit(3)

	response, err := c.Get("https://search.rakuten.co.jp/search/mall/%E3%83%8F%E3%83%BC%E3%83%90%E3%83%AA%E3%82%A6%E3%83%A0+%E3%82%AD%E3%83%83%E3%83%88/", nil)

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(response.String())

}
