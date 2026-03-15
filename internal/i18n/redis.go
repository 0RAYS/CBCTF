package i18n

var Redis = struct {
	// Key Error
	SetError string
	// Key Error
	GetError string
	// Key Error
	DeleteError string
}{
	SetError:    "redis.setError",
	GetError:    "redis.getError",
	DeleteError: "redis.deleteError",
}
