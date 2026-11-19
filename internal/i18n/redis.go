package i18n

var Redis = struct {
	SetError                 string
	GetError                 string
	NotFound                 string
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
