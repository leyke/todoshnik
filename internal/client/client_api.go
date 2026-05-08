package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
	"todoshnik/internal/bot/tg"
	apperrors "todoshnik/internal/errors"
)

type ApiClient struct {
	baseURL string
	client  *http.Client
}

func NewApiClient(baseURL string) *ApiClient {
	return &ApiClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *ApiClient) Get(ctx context.Context, endpoint string, query url.Values) (*http.Response, error) {
	url := c.baseURL + endpoint

	if query != nil {
		url += "?" + query.Encode()
	}
	fmt.Println(url)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// сервисный токен бота для внутренней проверки
	req.Header.Set("X-Bot-Service-Token", os.Getenv("BOT_SERVICE_TOKEN"))

	if token, ok := ctx.Value(tg.TokenContextKey).(string); ok {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, apperrors.ErrUnAuth
		}

		respBody, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf(
			"api error: %d: %s",
			resp.StatusCode,
			string(respBody),
		)
	}

	return resp, nil
}

func (c *ApiClient) Post(ctx context.Context, endpoint string, payload any) (*http.Response, error) {
	url := c.baseURL + endpoint
	fmt.Println(url)

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return nil, err
	}

	// сервисный токен бота для внутренней проверки
	req.Header.Set("X-Bot-Service-Token", os.Getenv("BOT_SERVICE_TOKEN"))

	req.Header.Set("Content-Type", "application/json")

	if token, ok := ctx.Value(tg.TokenContextKey).(string); ok {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, apperrors.ErrUnAuth
		}

		if resp.StatusCode == http.StatusUnauthorized {
			return nil, apperrors.ErrUnAuth
		}

		respBody, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf(
			"api error: %d: %s",
			resp.StatusCode,
			string(respBody),
		)
	}

	return resp, nil
}

func (c *ApiClient) Put(ctx context.Context, endpoint string, payload any) (*http.Response, error) {
	url := c.baseURL + endpoint
	fmt.Println(url)

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		url,
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return nil, err
	}

	// сервисный токен бота для внутренней проверки
	req.Header.Set("X-Bot-Service-Token", os.Getenv("BOT_SERVICE_TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	if token, ok := ctx.Value(tg.TokenContextKey).(string); ok {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, apperrors.ErrUnAuth
		}

		respBody, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf(
			"api error: %d: %s",
			resp.StatusCode,
			string(respBody),
		)
	}

	return resp, nil
}

func (c *ApiClient) Delete(ctx context.Context, endpoint string) (*http.Response, error) {
	url := c.baseURL + endpoint
	fmt.Println(url)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// сервисный токен бота для внутренней проверки
	req.Header.Set("X-Bot-Service-Token", os.Getenv("BOT_SERVICE_TOKEN"))

	if token, ok := ctx.Value(tg.TokenContextKey).(string); ok {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, apperrors.ErrUnAuth
		}

		respBody, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf(
			"api error: %d: %s",
			resp.StatusCode,
			string(respBody),
		)
	}

	return resp, nil
}
