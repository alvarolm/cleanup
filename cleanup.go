// Package cleanup provides utilities for managing deferred cleanup tasks
// with conditional execution based on error states.
package cleanup

// CleanTask represents a cleanup function that takes no parameters and returns nothing.
// These functions are executed during the cleanup phase based on error conditions.
type CleanTask func()

// Cleaner manages cleanup tasks that execute conditionally based on error state.
// Tasks are executed in LIFO (last-in-first-out) order when Clean() is called.
// This type is not safe for concurrent use.
type Cleaner struct {
	cleaned bool
	onerror []CleanTask
	onnil   []CleanTask
	always  []CleanTask
	Errptr  *error
}

// GetError returns the current error value from the Cleaner's error pointer.
// Returns nil if Errptr is nil or points to nil.
func (c *Cleaner) GetError() error {
	if c.Errptr != nil {
		return *c.Errptr
	}
	return nil
}

// ComplementOrSetErr updates the error pointer using a complementor function.
// If an error already exists, it calls complementor(existingErr, complementErr) to combine them.
// If no error exists, it sets complementErr as the error.
// Does nothing if Errptr is nil.
func (c *Cleaner) ComplementOrSetErr(complementor func(mainErr error, complementErr ...error) error, complementErr error) {
	if c.Errptr == nil {
		return
	}

	ierr := c.GetError()
	if ierr != nil {

		c.SetReturnError(complementor(ierr, complementErr))

	} else {

		c.SetReturnError(complementErr)

	}

}

// SetReturnError sets the error value in the Cleaner's error pointer.
// Does nothing if Errptr is nil.
func (c *Cleaner) SetReturnError(err error) {
	if c.Errptr != nil {
		*c.Errptr = err
	}
}

// NewCleaner creates a cleanup instance,
// use it at the beginning of your function
// hooking up your error pointer.
func NewCleaner(errPtr *error) (c *Cleaner) {
	c = new(Cleaner)
	c.Errptr = errPtr
	return
}

// Clean performs the cleanup by executing all registered tasks.
// Call it as a defer at the beginning of your function
// or within the function context.
// Tasks execute in LIFO (last-in-first-out) order.
// It executes only once; subsequent calls do nothing.
// If a task panics, the panic is recovered and remaining tasks continue to execute.
// This method is not safe for concurrent use.
func (c *Cleaner) Clean() {
	if c.cleaned {
		return
	}
	c.cleaned = true

	for i := len(c.always) - 1; i >= 0; i-- {
		func() {
			defer func() {
				recover()
			}()
			(c.always[i])()
		}()
	}
	if c.Errptr != nil {
		if (*c.Errptr) == nil {
			for i := len(c.onnil) - 1; i >= 0; i-- {
				func() {
					defer func() {
						recover()
					}()
					(c.onnil[i])()
				}()
			}
		} else {
			for i := len(c.onerror) - 1; i >= 0; i-- {
				func() {
					defer func() {
						recover()
					}()
					(c.onerror[i])()
				}()
			}
		}
	}
}

// OnError executes fx when err is not nil,
// just call it when you need it.
func (c *Cleaner) OnError(fx CleanTask) {
	c.onerror = append(c.onerror, fx)
}

// OnNil executes fx when err is nil,
// just call it when you need it.
func (c *Cleaner) OnNil(fx CleanTask) {
	c.onnil = append(c.onnil, fx)
}

// Always executes always,
// just call it when you need it.
func (c *Cleaner) Always(fx CleanTask) {
	c.always = append(c.always, fx)
}

// Always executes the provided function with the current error value.
// The function is called regardless of whether the error is nil or not.
// Does nothing if errPtr is nil.
func Always(errPtr *error, do func(error)) {
	if errPtr != nil {
		do(*errPtr)
	}
}

// OnError executes the provided function with the error value only if an error exists.
// Does nothing if errPtr is nil or points to nil.
func OnError(errPtr *error, do func(error)) {
	if errPtr != nil && (*errPtr) != nil {
		do(*errPtr)
	}
}

// OnNil executes the provided function only if no error exists.
// Does nothing if errPtr is nil or points to a non-nil error.
func OnNil(errPtr *error, do func()) {
	if errPtr != nil && (*errPtr) == nil {
		do()
	}
}

// OnTrue executes the provided function only if the boolean value is true.
// Does nothing if boolPtr is nil or points to false.
func OnTrue(boolPtr *bool, do func()) {
	if boolPtr != nil && *boolPtr {
		do()
	}
}

// OnFalse executes the provided function only if the boolean value is false.
// Does nothing if boolPtr is nil or points to true.
func OnFalse(boolPtr *bool, do func()) {
	if boolPtr != nil && !(*boolPtr) {
		do()
	}
}

// ExecIfSet executes the function if both the pointer and the function it points to are non-nil.
// Does nothing if funcPtr is nil or points to a nil function.
func ExecIfSet(funcPtr *func()) {
	if funcPtr != nil && *funcPtr != nil {
		(*funcPtr)()
	}
}
