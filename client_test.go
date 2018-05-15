package luminati

import (
	"fmt"
	"testing"
)

func TestLuminatiClient_ShoulReturnContent_When200(t *testing.T) {
	c := NewClient("lum-customer-xxx-zone-residential", "xxx")

	c.SetFailuresLimit(3)

	link := "https://baidu.com"
	link = "http://lumtest.com/myip.json"
	response, err := c.Get(link, nil)

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(response.String())

	response, err = c.Get(link, nil)

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(response.String())
	for k, v := range response.RawResponse.Request.Header {
		fmt.Printf("%s: %v\n", k, v)
	}

}
