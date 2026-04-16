package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/ctxutil"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/logger"
)

type UsersClient struct {
	baseURL string
	client  *http.Client
}

func NewUsersClient(baseURL string) *UsersClient {
	return &UsersClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

// implementing interface UserValidator
func (c *UsersClient) ValidateUserExists(ctx context.Context, userID string) error {
	url := fmt.Sprintf("%s/%s", c.baseURL, userID)

	maxRetries := 3
	retryDelay := 100 * time.Millisecond

	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		corrID := ctxutil.GetCorrelationID(ctx)
		req.Header.Set("X-Correlation-Id", corrID)

		resp, err := c.client.Do(req)

		if err != nil {
			lastErr = err
			logger.Printf(ctx, "[UsersClient] Attempt %d failed: %v", attempt+1, err)

			if attempt < maxRetries {
				time.Sleep(retryDelay)
				retryDelay *= 2
			}
			continue
		}

		if resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}

		if resp.StatusCode == http.StatusNotFound {
			resp.Body.Close()
			return shared.NewValidationError("owner_user_id '%s' does not exist", userID)
		}

		if resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("service returned transient status: %d", resp.StatusCode)
			logger.Printf(ctx, "[UsersClient] Attempt %d failed: %v", attempt+1, lastErr)

			resp.Body.Close()
			if attempt < maxRetries {
				time.Sleep(retryDelay)
				retryDelay *= 2
			}
			continue
		}

		status := resp.StatusCode
		resp.Body.Close()
		return shared.NewServiceUnavailableError("users service returned unexpected status: %d", status)
	}

	var netErr net.Error
	if errors.Is(lastErr, context.DeadlineExceeded) || (errors.As(lastErr, &netErr) && netErr.Timeout()) {
		return shared.NewGatewayTimeoutError("users service timed out after %d retries", maxRetries)
	}

	return shared.NewServiceUnavailableError("users service is unreachable after %d retries: %v", maxRetries, lastErr)
}
