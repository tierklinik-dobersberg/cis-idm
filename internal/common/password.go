package common

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func (c *Service) GeneratePasswordResetToken(ctx context.Context, userID string) (string, string, error) {
	// Generate a new security code
	source := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(source)
	code := fmt.Sprintf("%d", rand.Intn(999999-100000)+100000)

	// store the code in cache
	cacheKey := fmt.Sprintf("password-reset:%s", code)

	if err := c.cache.PutKeyTTL(ctx, cacheKey, userID, time.Hour*24); err != nil {
		return "", "", err
	}

	return code, cacheKey, nil
}
