package redis

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	ErrGeneratorPoolEmpty   = errors.New("generator pool is empty")
	ErrNoAvailableGenerator = errors.New("no available generator")
)

const (
	GeneratorKey               = "generator:%d"
	GeneratorAttachmentLockKey = "generator:%d:locked"
	generatorChallengeSetKey   = "generators:challenge:%d"
	generatorContestSetKey     = "generators:contest:%d:challenge:%d"
	generatorLockTTL           = 5 * time.Minute
)

const lockAvailableGeneratorScript = `
local ids = redis.call('SMEMBERS', KEYS[1])
local count = #ids
if count == 0 then
    return {-1, ''}
end

local start = math.random(1, count)
for i = 0, count - 1 do
    local index = ((start + i - 1) % count) + 1
    local id = ids[index]
    local generator_key = 'generator:' .. id
    local generator = redis.call('GET', generator_key)
    if generator then
        local lock_key = generator_key .. ':locked'
        if redis.call('SET', lock_key, ARGV[1], 'NX', 'PX', ARGV[2]) then
            return {1, generator}
        end
    else
        redis.call('SREM', KEYS[1], id)
    end
end

return {0, ''}
`

const unlockGeneratorAttachmentScript = `
if redis.call('GET', KEYS[1]) == ARGV[1] then
    return redis.call('DEL', KEYS[1])
end
return 0
`

const refreshGeneratorAttachmentScript = `
if redis.call('GET', KEYS[1]) == ARGV[1] then
    return redis.call('PEXPIRE', KEYS[1], ARGV[2])
end
return 0
`

func RegisterGenerator(ctx context.Context, generator model.Generator) error {
	data, err := json.Marshal(generator)
	if err != nil {
		return fmt.Errorf("marshal generator failed: %w", err)
	}
	pipe := RDB.TxPipeline()
	pipe.Set(ctx, fmt.Sprintf(GeneratorKey, generator.ID), data, 0)
	pipe.SAdd(ctx, generatorSetKey(generator.ContestID.V, generator.ContestID.Valid, generator.ChallengeID), generator.ID)
	if _, err = pipe.Exec(ctx); err != nil {
		log.Logger.Warningf("Failed to register generator in redis: generator_id=%d err=%v", generator.ID, err)
		return err
	}
	return nil
}

func UnregisterGenerator(ctx context.Context, generator model.Generator) error {
	pipe := RDB.TxPipeline()
	pipe.Del(ctx, fmt.Sprintf(GeneratorKey, generator.ID))
	pipe.Del(ctx, fmt.Sprintf(GeneratorAttachmentLockKey, generator.ID))
	pipe.SRem(ctx, generatorSetKey(generator.ContestID.V, generator.ContestID.Valid, generator.ChallengeID), generator.ID)
	if _, err := pipe.Exec(ctx); err != nil {
		log.Logger.Warningf("Failed to unregister generator in redis: generator_id=%d err=%v", generator.ID, err)
		return err
	}
	return nil
}

func LockAvailableGenerator(ctx context.Context, contestID, challengeID uint) (model.Generator, string, error) {
	token := fmt.Sprintf("%d:%d:%d", contestID, challengeID, time.Now().UnixNano())
	result, err := RDB.Eval(
		ctx,
		lockAvailableGeneratorScript,
		[]string{generatorSetKey(contestID, contestID > 0, challengeID)},
		token,
		int64(generatorLockTTL/time.Millisecond),
	).Result()
	if err != nil {
		log.Logger.Warningf("Failed to lock available generator: contest_id=%d challenge_id=%d err=%v", contestID, challengeID, err)
		return model.Generator{}, "", err
	}

	resultSlice, ok := result.([]interface{})
	if !ok || len(resultSlice) != 2 {
		return model.Generator{}, "", fmt.Errorf("invalid lock generator script result")
	}
	status, ok := resultSlice[0].(int64)
	if !ok {
		return model.Generator{}, "", fmt.Errorf("invalid lock generator status")
	}
	switch status {
	case -1:
		return model.Generator{}, "", ErrGeneratorPoolEmpty
	case 0:
		return model.Generator{}, "", ErrNoAvailableGenerator
	}

	data, ok := resultSlice[1].(string)
	if !ok || data == "" {
		return model.Generator{}, "", fmt.Errorf("invalid lock generator payload")
	}
	var generator model.Generator
	if err = json.Unmarshal([]byte(data), &generator); err != nil {
		return model.Generator{}, "", fmt.Errorf("unmarshal generator failed: %w", err)
	}
	return generator, token, nil
}

func LockGeneratorAttachment(ctx context.Context, generatorID uint) (string, error) {
	key := fmt.Sprintf(GeneratorAttachmentLockKey, generatorID)
	token := fmt.Sprintf("%d:%d", generatorID, time.Now().UnixNano())
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		ok, err := RDB.SetNX(ctx, key, token, generatorLockTTL).Result()
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

func RefreshGeneratorAttachmentLock(ctx context.Context, generatorID uint, token string) (bool, error) {
	key := fmt.Sprintf(GeneratorAttachmentLockKey, generatorID)
	result, err := RDB.Eval(ctx, refreshGeneratorAttachmentScript, []string{key}, token, int64(generatorLockTTL/time.Millisecond)).Int()
	if err != nil {
		log.Logger.Warningf("Failed to refresh generator attachment lock: key=%s err=%v", key, err)
		return false, err
	}
	return result == 1, nil
}

func UnlockGeneratorAttachment(ctx context.Context, generatorID uint, token string) error {
	key := fmt.Sprintf(GeneratorAttachmentLockKey, generatorID)
	if _, err := RDB.Eval(ctx, unlockGeneratorAttachmentScript, []string{key}, token).Result(); err != nil {
		log.Logger.Warningf("Failed to unlock generator attachment: key=%s err=%v", key, err)
		return err
	}
	return nil
}

func generatorSetKey(contestID uint, contestValid bool, challengeID uint) string {
	if contestValid && contestID > 0 {
		return fmt.Sprintf(generatorContestSetKey, contestID, challengeID)
	}
	return fmt.Sprintf(generatorChallengeSetKey, challengeID)
}
