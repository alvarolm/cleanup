package cleanup

import "testing"

// TestNilPointerSafety verifies that functions handle nil pointers safely
func TestNilPointerSafety(t *testing.T) {
	// Test SetReturnError with nil Errptr
	c := &Cleaner{Errptr: nil}
	c.SetReturnError(nil) // Should not panic

	// Test ComplementOrSetErr with nil Errptr
	c.ComplementOrSetErr(func(mainErr error, complementErr ...error) error {
		return mainErr
	}, nil) // Should not panic

	// Test Always with nil errPtr
	Always(nil, func(error) {}) // Should not panic

	// Test OnError with nil errPtr
	OnError(nil, func(error) {}) // Should not panic

	// Test OnNil with nil errPtr
	OnNil(nil, func() {}) // Should not panic

	// Test OnTrue with nil boolPtr
	OnTrue(nil, func() {}) // Should not panic

	// Test OnFalse with nil boolPtr
	OnFalse(nil, func() {}) // Should not panic

	// Test ExecIfSet with nil funcPtr
	ExecIfSet(nil) // Should not panic

	// Test ExecIfSet with non-nil pointer to nil function
	var nilFunc func()
	ExecIfSet(&nilFunc) // Should not panic

	t.Log("All nil pointer safety tests passed")
}
