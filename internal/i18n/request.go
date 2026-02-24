package i18n

var Response = struct {
	// Error
	BadRequest      string
	Unauthorized    string
	Forbidden       string
	TooManyRequests string
}{
	BadRequest:      "response.badRequest",
	Unauthorized:    "response.unauthorized",
	Forbidden:       "response.forbidden",
	TooManyRequests: "response.tooManyRequests",
}
