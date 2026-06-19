package redis

import (
	"CBCTF/internal/log"
	"context"
	"fmt"
	"time"
)

const GeneratorAttachmentLockKey = "generator:%d:locked"

const unlockGeneratorAttachmentScript = `
if redis.call('GET', KEYS[1]) == ARGV[1] then
    return redis.call('DEL', KEYS[1])
end
return 0
`

func LockGeneratorAttachment(ctx context.Context, generatorID uint) (string, error) {
	key := fmt.Sprintf(GeneratorAttachmentLockKey, generatorID)
	token := fmt.Sprintf("%d:%d", generatorID, time.Now().UnixNano())
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		ok, err := RDB.SetNX(ctx, key, token, 3*time.Minute).Result()
		if err != nil {
			log.Logger.Warningf("Failed to lock generator attachment: key=%s err=%v", key, err)
			return "", err
		}
		if ok {
			return token, nil
		}

		select {
		case <-ctx.Done():
			return "", fmt.Errorf("lock generator attachment timed out: %w", ctx.Err())
		case <-ticker.C:
		}
	}
}

func UnlockGeneratorAttachment(ctx context.Context, generatorID uint, token string) error {
	key := fmt.Sprintf(GeneratorAttachmentLockKey, generatorID)
	if _, err := RDB.Eval(ctx, unlockGeneratorAttachmentScript, []string{key}, token).Result(); err != nil {
		log.Logger.Warningf("Failed to unlock generator attachment: key=%s err=%v", key, err)
		return err
	}
	return nil
}
