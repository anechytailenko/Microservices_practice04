package infrastructure

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"
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
	url := fmt.Sprintf("%s/users/%s", c.baseURL, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return shared.NewServiceUnavailableError("users service is unreachable: %v", err)
	}

	defer resp.Body.Close()

	//  404 -> 400
	if resp.StatusCode == http.StatusNotFound {
		return shared.NewValidationError("owner_user_id '%s' does not exist", userID)
	}
	if resp.StatusCode != http.StatusOK {
		return shared.NewServiceUnavailableError("users service returned unexpected status: %d", resp.StatusCode)
	}

	return nil
}
