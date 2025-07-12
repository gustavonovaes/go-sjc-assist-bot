package news

import "time"

const (
	// HTTP client
	defaultRequestTimeout = 10 * time.Second

	// Meon news
	meonURL    = "https://www.meon.com.br/noticias/rmvale/rss"
	meonOrigin = "Meon"

	// Sampi news
	sampiURL          = "https://sampi.net.br/ovale/categoria/ultimas"
	sampiOrigin       = "Sampi"
	sampiRegexPattern = `(?s)<div class="notia\s*">\s*<a[^>]*href="([^"]*)"[^>]*>.*?<h3[^>]*>(?:<span[^>]*>[^<]*</span>)?(.*?)(?:<time[^>]*>[^<]*</time>)?</h3>`
)
