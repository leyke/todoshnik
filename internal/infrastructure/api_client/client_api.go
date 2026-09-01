package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	authcontext "todoshnik/internal/infrastructure/context_manager/auth"
)

const (
	headerServiceToken string = "X-Bot-Service-Token"
	headerAuth         string = "Authorization"
	headerContentType  string = "Content-Type"
)

type ApiClient struct {
	baseURL      string
	serviceToken string
	client       *http.Client
}

func NewApiClient(baseURL string, serviceToken string, requestTimeout time.Duration) *ApiClient {
	return &ApiClient{
		baseURL:      baseURL,
		serviceToken: serviceToken,
		client: &http.Client{
			Timeout: requestTimeout,
		},
	}
}

func (c *ApiClient) Get(ctx context.Context, endpoint string, query url.Values) (*http.Response, error) {
	url := c.baseURL + endpoint

	if query != nil {
		url += "?" + query.Encode()
	}

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
	req.Header.Set(headerServiceToken, c.serviceToken)
	if token, ok := authcontext.GetToken(ctx); ok {
		req.Header.Set(headerAuth, "Bearer "+token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, ErrUnAuth
		}

		respBody, _ := io.ReadAll(resp.Body)
		defer func() {
			_ = resp.Body.Close()
		}()

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
	req.Header.Set(headerServiceToken, c.serviceToken)
	req.Header.Set(headerContentType, "application/json")

	if token, ok := authcontext.GetToken(ctx); ok {
		req.Header.Set(headerAuth, "Bearer "+token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, ErrUnAuth
		}

		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}

		respBody, _ := io.ReadAll(resp.Body)
		defer func() {
			_ = resp.Body.Close()
		}()

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
	req.Header.Set(headerServiceToken, c.serviceToken)
	req.Header.Set(headerContentType, "application/json")

	if token, ok := authcontext.GetToken(ctx); ok {
		req.Header.Set(headerAuth, "Bearer "+token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, ErrUnAuth
		}

		respBody, _ := io.ReadAll(resp.Body)
		defer func() {
			_ = resp.Body.Close()
		}()

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
	req.Header.Set(headerServiceToken, c.serviceToken)
	if token, ok := authcontext.GetToken(ctx); ok {
		req.Header.Set(headerAuth, "Bearer "+token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, ErrUnAuth
		}

		respBody, _ := io.ReadAll(resp.Body)
		defer func() {
			_ = resp.Body.Close()
		}()

		return nil, fmt.Errorf(
			"api error: %d: %s",
			resp.StatusCode,
			string(respBody),
		)
	}

	return resp, nil
}
