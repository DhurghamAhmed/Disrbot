package utils

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/chromedp/chromedp"
)

func GenerateCarbonImage(code string, theme string) ([]byte, error) {
	params := url.Values{}
	params.Add("code", code)
	params.Add("t", theme)
	targetURL := fmt.Sprintf("https://carbon.now.sh/?%s", params.Encode())

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var buf []byte

	err := chromedp.Run(ctx,
		chromedp.Navigate(targetURL),
		chromedp.WaitVisible("#export-container", chromedp.ByID),
		chromedp.Screenshot("#export-container", &buf, chromedp.NodeVisible, chromedp.ByID),
	)

	if err != nil {
		return nil, err
	}

	return buf, nil
}