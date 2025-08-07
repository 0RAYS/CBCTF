package redis

import (
	"context"
	"fmt"
	"time"
)

func SetOauthState(provider, state, verifier string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return RDB.Set(ctx, fmt.Sprintf("oauth:%s:%s", provider, state), verifier, 10*time.Minute).Err()
}

func GetOauthVerifier(provider, state string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return RDB.Get(ctx, fmt.Sprintf("oauth:%s:%s", provider, state)).Result()
}

func DelOauthState(provider string, state string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return RDB.SRem(ctx, fmt.Sprintf("oauth:%s", provider), state).Err()
}
