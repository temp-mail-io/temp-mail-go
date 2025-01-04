package temp_mail_go

import (
	"context"
	"fmt"
)

func ExampleClient_RateLimit() {
	c := NewClient("YOUR_API_KEY", nil)
	resp, err := c.RateLimit(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("limit: %d, used: %d, remaining: %d, reset: %d", resp.Limit, resp.Used, resp.Remaining, resp.Reset)
	// Output: limit: 1000, used: 0, remaining: 1000, reset: 1738367999
}
