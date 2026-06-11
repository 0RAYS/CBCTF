package i18n

var Redis = struct {
	// Key Error
	SetError string
	// Key Error
	GetError string
	// Key
	NotFound string
	// Key Error
	DeleteError              string
	InvalidScriptResult      string
	InvalidScriptPort        string
	InvalidScriptSuccessFlag string
	NoAvailablePort          string
}{
	SetError:                 "redis.setError",
	GetError:                 "redis.getError",
	NotFound:                 "redis.notFound",
	DeleteError:              "redis.deleteError",
	InvalidScriptResult:      "redis.invalidScriptResult",
	InvalidScriptPort:        "redis.invalidScriptPort",
	InvalidScriptSuccessFlag: "redis.invalidScriptSuccessFlag",
	NoAvailablePort:          "redis.noAvailablePort",
}
