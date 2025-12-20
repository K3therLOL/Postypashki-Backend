package curl

import (
	"context"
	"fmt"
	"hedgedcurl/fail"
	"hedgedcurl/timeout"
	"hedgedcurl/usage"
	"io"
	"log"
	"os"
	"net/http"
	"strings"
	"sync"
	"time"
)

func Get(ctx context.Context, url string) *http.Response {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil 
	}

	return resp
}

func formatResponse(resp *http.Response) {
	fmt.Println(resp.Status)
	fmt.Println(resp.Proto)
	for key, values := range resp.Header {
		formatValue := strings.Join(values, ", ")
		fmt.Printf("%s: %s\n", key, formatValue)
	}
	fmt.Println()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()
	fmt.Println(string(body))
}

func hedged(urls []string) error {
	timeout := time.Duration(usage.Time) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Go(func() {
			resp := Get(ctx, url)
			if resp != nil {
				formatResponse(resp)
				cancel()
			}
		})
	 }

	 wg.Wait()
	 return ctx.Err()
}

func FormatHedged(urls []string) int {
	switch err := hedged(urls); err {
		case context.DeadlineExceeded:
			fmt.Fprintln(os.Stderr, timeout.Error())
			return timeout.RequestsTimeout
		case nil:
			fmt.Fprintln(os.Stderr, fail.Error())
			return http.StatusBadRequest
	}

	return 0
}
