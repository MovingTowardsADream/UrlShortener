package usecase

import "time"

type Option func(*UrlsUseCase)

func Timeout(timeout time.Duration) Option {
	return func(uc *UrlsUseCase) {
		uc.timeout = timeout
	}
}
