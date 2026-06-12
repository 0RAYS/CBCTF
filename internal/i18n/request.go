package i18n

var Response = struct {
	// Error
	BadRequest      string
	Unauthorized    string
	Forbidden       string
	RequestTooLarge string
	TooManyRequests string
}{
	BadRequest:      "response.badRequest",
	Unauthorized:    "response.unauthorized",
	Forbidden:       "response.forbidden",
	RequestTooLarge: "response.requestTooLarge",
	TooManyRequests: "response.tooManyRequests",
}
