package resp

// Keep the base object visible while distinguishing unavailable optional data
// from a database-confirmed zero/empty value.
func markUnavailable(data map[string]any, fields []string) {
	if len(fields) == 0 {
		return
	}
	data["unavailable"] = fields
	for _, field := range fields {
		data[field] = nil
	}
}
