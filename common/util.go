package common

func SliceToChannel[T any](tokens []T) chan T {
	tokenChan := make(chan T)
	go func() {
		for _, token := range tokens {
			tokenChan <- token
		}
		close(tokenChan)
	}()
	return tokenChan
}
