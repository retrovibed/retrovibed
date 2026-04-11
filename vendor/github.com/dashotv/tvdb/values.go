package tvdb

// StringValue returns the value of pointer if set, otherwise the zero value
func StringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

// BoolValue returns the value of pointer if set, otherwise the zero value
func BoolValue(v *bool) bool {
	if v == nil {
		return false
	}
	return *v
}

// IntValue returns the value of pointer if set, otherwise the zero value
func IntValue(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

// Int64Value returns the value of pointer if set, otherwise the zero value
func Int64Value(v *int64) int64 {
	if v == nil {
		return 0
	}
	return *v
}

// Float32Value returns the value of pointer if set, otherwise the zero value
func Float32Value(v *float32) float32 {
	if v == nil {
		return 0
	}
	return *v
}

// Float64Value returns the value of pointer if set, otherwise the zero value
func Float64Value(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}
