package goclient

import (
	"fmt"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	c := NewTimeoutClient(time.Second)
	c.SetHeader("Host", "github.com")
	u := "http://127.0.0.1:1111"
	fmt.Println(c.GetData(u))
}
