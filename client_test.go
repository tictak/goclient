package dialog

import (
	"fmt"
	"net/url"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	c := NewTimeoutClient(time.Second)
	c.SetHeader("Host", "github.com")
	u, e := url.Parse("http://127.0.0.1:1111")
	fmt.Println(e)
	fmt.Println(c.GetData(u))
}
