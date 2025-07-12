package news

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type SampiNewsClient struct {
	httpClient *http.Client
	newsRegex  *regexp.Regexp
}

func NewSampiNewsClient() *SampiNewsClient {
	return &SampiNewsClient{
		httpClient: &http.Client{Timeout: defaultRequestTimeout},
		newsRegex:  regexp.MustCompile(sampiRegexPattern),
	}
}

func GetSampiNews() ([]News, error) {
	client := NewSampiNewsClient()
	return client.FetchNews(context.Background())
}

func (c *SampiNewsClient) FetchNews(ctx context.Context) ([]News, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sampiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch news from %s: %w", sampiURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d from %s", resp.StatusCode, sampiURL)
	}

	htmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	news, err := c.parseHTML(htmlData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	if len(news) == 0 {
		return nil, fmt.Errorf("no news items found in HTML")
	}

	return news, nil
}

func (c *SampiNewsClient) parseHTML(htmlData []byte) ([]News, error) {
	matches := c.newsRegex.FindAllSubmatch(htmlData, -1)

	if matches == nil {
		return nil, fmt.Errorf("no news items found in HTML")
	}

	news := make([]News, 0, len(matches))
	for _, match := range matches {
		if len(match) < 3 {
			log.Printf("Skipping malformed news item match: %v", match)
			continue
		}

		link := strings.TrimSpace(string(match[1]))
		title := stripHtmlTags(string(match[2]))
		if link == "" || title == "" {
			log.Printf("Skipping news item with empty link or title: %s", title)
			continue
		}

		newsItem := News{
			Link:    link,
			Title:   title,
			Content: "",
			Origin:  sampiOrigin,
			Tags:    []string{},
		}

		news = append(news, newsItem)
	}

	return news, nil
}
