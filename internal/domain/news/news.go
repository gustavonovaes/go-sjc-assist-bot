package news

import (
	"fmt"
	"sync"
)

type NewsProvider struct {
	GetNews func() ([]News, error)
}

func GetLastNews() ([]News, error) {
	providers := []NewsProvider{
		{GetNews: GetMeonNews},
		{GetNews: GetSampiNews},
	}

	var (
		allNews []News
		errs    []error
		wg      sync.WaitGroup
		mu      sync.Mutex
	)

	wg.Add(len(providers))
	for _, provider := range providers {
		go func(p NewsProvider) {
			defer wg.Done()
			news, err := p.GetNews()

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				errs = append(errs, fmt.Errorf("failed to get news: %w", err))
				return
			}

			allNews = append(allNews, news...)
		}(provider)
	}

	wg.Wait()

	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to get news from providers: %v", errs)
	}

	return allNews, nil
}
