package utils

import (
	"errors"
	"io"
	"net/http"
	"net/url"
)

func GetShortUrl(originalURL string) (string, error) {
	apiURL := "https://clck.ru/--"
	resp, err := http.PostForm(apiURL, url.Values{"url": {originalURL}})
	if err != nil {
		return "", errors.New("ошибка api clck.ru")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(body), nil
	}

	return "", errors.New("ошибка api clck.ru")
}
