package session

import (
	"fmt"
	"log"
	"net/url"
	"testing"
	"time"
)

func TestSession(t *testing.T) {

	u, err := url.Parse("http://www.youkeda.com:3000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(GetDomain(u.Hostname()))
	// for i := 0; i < 100; i++ {
	// 	NewSession()
	// }

	fmt.Println(time.Unix(1, 0))
}
