package api

import (
	"net/http"
	"time"
)

func createHeaders(token string) http.Header {
	headers := http.Header{}
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:54.0) Gecko/20100101 Firefox/54.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	}
	headers.Set("User-Agent", userAgents[time.Now().UnixNano()%int64(len(userAgents))])
	headers.Set("accept", "application/json")
	headers.Set("accept-encoding", "gzip, deflate, br, zstd")
	headers.Set("accept-language", "Bid-ID,id;q=0.6")
	headers.Set("content-type", "application/json")
	headers.Set("origin", "https://www.magicnewton.com")
	headers.Set("referer", "https://www.magicnewton.com/portal/rewards")
	headers.Set("Cookie", "__Secure-next-auth.session-token="+token)
	return headers
}
