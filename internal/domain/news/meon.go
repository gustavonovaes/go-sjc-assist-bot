package news

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type MeonNewsClient struct {
	httpClient *http.Client
}

func NewMeonNewsClient() *MeonNewsClient {
	return &MeonNewsClient{
		httpClient: &http.Client{Timeout: defaultRequestTimeout},
	}
}

func GetMeonNews() ([]News, error) {
	client := NewMeonNewsClient()
	return client.FetchNews(context.Background())
}

func (c *MeonNewsClient) FetchNews(ctx context.Context) ([]News, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, meonURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/rss+xml, application/xml, text/xml")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch news from %s: %w", meonURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d from %s", resp.StatusCode, meonURL)
	}

	xmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	news, err := c.parseRSSXML(xmlData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS XML: %w", err)
	}

	if len(news) == 0 {
		return nil, fmt.Errorf("no news items found in RSS feed")
	}

	return news, nil
}

func (c *MeonNewsClient) parseRSSXML(xmlData []byte) ([]News, error) {
	var rss MeonRssChannel

	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	if err := decoder.Decode(&rss); err != nil {
		return nil, fmt.Errorf("failed to decode RSS XML: %w", err)
	}

	news := make([]News, 0, len(rss.Channel.Items))
	for _, item := range rss.Channel.Items {
		news = append(news, c.convertRSSItemToNews(item))
	}

	return news, nil
}

func (c *MeonNewsClient) convertRSSItemToNews(item MeonRssItem) News {
	return News{
		Title:   stripHtmlTags(item.Title),
		Link:    strings.TrimSpace(item.Link),
		Content: stripHtmlTags(item.Description),
		Origin:  meonOrigin,
		Tags:    []string{},
	}
}
