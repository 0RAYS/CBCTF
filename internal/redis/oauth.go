package redis

import (
	"context"
	"fmt"
	"time"
)

const oauthKey = "oauth:%s:%s"

func SetOauthState(provider, state, verifier string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return RDB.Set(ctx, fmt.Sprintf(oauthKey, provider, state), verifier, 10*time.Minute).Err()
}

func GetOauthVerifier(provider, state string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return RDB.Get(ctx, fmt.Sprintf(oauthKey, provider, state)).Result()
}

func DelOauthState(provider string, state string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return RDB.Del(ctx, fmt.Sprintf(oauthKey, provider, state)).Err()
}
