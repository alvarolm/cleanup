package cleanup

import (
	"errors"
	"testing"
)

// TestSuccessCase tests that OnNil is executed when no error occurs
func TestSuccessCase(t *testing.T) {
	var err error
	cleaner := NewCleaner(&err)

	executed := false
	cleaner.OnNil(func() {
		executed = true
	})

	cleaner.Clean()

	if !executed {
		t.Error("OnNil should have been executed")
	}
}

// TestErrorCase tests that OnError is executed when an error occurs
func TestErrorCase(t *testing.T) {
	var err error = errors.New("test error")
	cleaner := NewCleaner(&err)

	executed := false
	cleaner.OnError(func() {
		executed = true
	})

	cleaner.Clean()

	if !executed {
		t.Error("OnError should have been executed")
	}
}

// TestAlwaysExecutes tests that Always tasks execute regardless of error state
func TestAlwaysExecutes(t *testing.T) {
	t.Run("with error", func(t *testing.T) {
		var err error = errors.New("test error")
		cleaner := NewCleaner(&err)

		executed := false
		cleaner.Always(func() {
			executed = true
		})

		cleaner.Clean()

		if !executed {
			t.Error("Always should have been executed even with error")
		}
	})

	t.Run("without error", func(t *testing.T) {
		var err error
		cleaner := NewCleaner(&err)

		executed := false
		cleaner.Always(func() {
			executed = true
		})

		cleaner.Clean()

		if !executed {
			t.Error("Always should have been executed without error")
		}
	})
}

// TestLIFOExecutionOrder tests that tasks execute in LIFO order
func TestLIFOExecutionOrder(t *testing.T) {
	var err error
	cleaner := NewCleaner(&err)

	var order []int

	cleaner.Always(func() { order = append(order, 1) })
	cleaner.Always(func() { order = append(order, 2) })
	cleaner.Always(func() { order = append(order, 3) })

	cleaner.Clean()

	if len(order) != 3 {
		t.Fatalf("Expected 3 tasks to execute, got %d", len(order))
	}

	// LIFO: last added (3) executes first
	if order[0] != 3 || order[1] != 2 || order[2] != 1 {
		t.Errorf("Expected order [3, 2, 1], got %v", order)
	}
}

// TestMultipleCleanCalls tests that Clean() only executes once
func TestMultipleCleanCalls(t *testing.T) {
	var err error
	cleaner := NewCleaner(&err)

	count := 0
	cleaner.Always(func() {
		count++
	})

	cleaner.Clean()
	cleaner.Clean()
	cleaner.Clean()

	if count != 1 {
		t.Errorf("Expected Clean to execute tasks once, executed %d times", count)
	}
}

// TestPanicRecovery tests that panics in tasks are recovered
func TestPanicRecovery(t *testing.T) {
	var err error
	cleaner := NewCleaner(&err)

	var executed []int

	cleaner.Always(func() { executed = append(executed, 1) })
	cleaner.Always(func() { panic("test panic") })
	cleaner.Always(func() { executed = append(executed, 3) })

	cleaner.Clean()

	// LIFO order: 3, panic, 1
	// All should execute despite panic in middle
	if len(executed) != 2 {
		t.Errorf("Expected 2 tasks to execute (panic recovered), got %d", len(executed))
	}
	if len(executed) == 2 && (executed[0] != 3 || executed[1] != 1) {
		t.Errorf("Expected [3, 1], got %v", executed)
	}
}

// TestGetError tests the GetError method
func TestGetError(t *testing.T) {
	t.Run("with error", func(t *testing.T) {
		testErr := errors.New("test error")
		err := testErr
		cleaner := NewCleaner(&err)

		if cleaner.GetError() != testErr {
			t.Error("GetError should return the error")
		}
	})

	t.Run("without error", func(t *testing.T) {
		var err error
		cleaner := NewCleaner(&err)

		if cleaner.GetError() != nil {
			t.Error("GetError should return nil when no error")
		}
	})

	t.Run("with nil errptr", func(t *testing.T) {
		cleaner := &Cleaner{Errptr: nil}

		if cleaner.GetError() != nil {
			t.Error("GetError should return nil when Errptr is nil")
		}
	})
}

// TestSetReturnError tests the SetReturnError method
func TestSetReturnError(t *testing.T) {
	t.Run("sets error", func(t *testing.T) {
		var err error
		cleaner := NewCleaner(&err)

		testErr := errors.New("test error")
		cleaner.SetReturnError(testErr)

		if err != testErr {
			t.Error("SetReturnError should set the error")
		}
	})

	t.Run("with nil errptr", func(t *testing.T) {
		cleaner := &Cleaner{Errptr: nil}

		// Should not panic
		cleaner.SetReturnError(errors.New("test"))
	})
}

// TestComplementOrSetErr tests the ComplementOrSetErr method
func TestComplementOrSetErr(t *testing.T) {
	complementor := func(mainErr error, complementErr ...error) error {
		if len(complementErr) > 0 {
			return errors.New(mainErr.Error() + ": " + complementErr[0].Error())
		}
		return mainErr
	}

	t.Run("with existing error", func(t *testing.T) {
		var err error = errors.New("main error")
		cleaner := NewCleaner(&err)

		cleaner.ComplementOrSetErr(complementor, errors.New("complement"))

		if err.Error() != "main error: complement" {
			t.Errorf("Expected 'main error: complement', got '%s'", err.Error())
		}
	})

	t.Run("without existing error", func(t *testing.T) {
		var err error
		cleaner := NewCleaner(&err)

		complementErr := errors.New("new error")
		cleaner.ComplementOrSetErr(complementor, complementErr)

		if err != complementErr {
			t.Error("Should set the complement error when no main error")
		}
	})

	t.Run("with nil errptr", func(t *testing.T) {
		cleaner := &Cleaner{Errptr: nil}

		// Should not panic
		cleaner.ComplementOrSetErr(complementor, errors.New("test"))
	})
}

// TestHelperFunctions tests the standalone helper functions
func TestAlwaysHelper(t *testing.T) {
	t.Run("executes with error", func(t *testing.T) {
		err := errors.New("test")
		executed := false
		Always(&err, func(e error) {
			executed = true
			if e != err {
				t.Error("Should receive the error")
			}
		})
		if !executed {
			t.Error("Should have executed")
		}
	})

	t.Run("executes with nil error", func(t *testing.T) {
		var err error
		executed := false
		Always(&err, func(e error) {
			executed = true
			if e != nil {
				t.Error("Should receive nil")
			}
		})
		if !executed {
			t.Error("Should have executed")
		}
	})

	t.Run("safe with nil pointer", func(t *testing.T) {
		Always(nil, func(e error) {
			t.Error("Should not execute")
		})
	})
}

func TestOnErrorHelper(t *testing.T) {
	t.Run("executes with error", func(t *testing.T) {
		err := errors.New("test")
		executed := false
		OnError(&err, func(e error) {
			executed = true
		})
		if !executed {
			t.Error("Should have executed")
		}
	})

	t.Run("does not execute without error", func(t *testing.T) {
		var err error
		OnError(&err, func(e error) {
			t.Error("Should not execute")
		})
	})

	t.Run("safe with nil pointer", func(t *testing.T) {
		OnError(nil, func(e error) {
			t.Error("Should not execute")
		})
	})
}

func TestOnNilHelper(t *testing.T) {
	t.Run("executes without error", func(t *testing.T) {
		var err error
		executed := false
		OnNil(&err, func() {
			executed = true
		})
		if !executed {
			t.Error("Should have executed")
		}
	})

	t.Run("does not execute with error", func(t *testing.T) {
		err := errors.New("test")
		OnNil(&err, func() {
			t.Error("Should not execute")
		})
	})

	t.Run("safe with nil pointer", func(t *testing.T) {
		OnNil(nil, func() {
			t.Error("Should not execute")
		})
	})
}

func TestOnTrueHelper(t *testing.T) {
	t.Run("executes with true", func(t *testing.T) {
		val := true
		executed := false
		OnTrue(&val, func() {
			executed = true
		})
		if !executed {
			t.Error("Should have executed")
		}
	})

	t.Run("does not execute with false", func(t *testing.T) {
		val := false
		OnTrue(&val, func() {
			t.Error("Should not execute")
		})
	})

	t.Run("safe with nil pointer", func(t *testing.T) {
		OnTrue(nil, func() {
			t.Error("Should not execute")
		})
	})
}

func TestOnFalseHelper(t *testing.T) {
	t.Run("executes with false", func(t *testing.T) {
		val := false
		executed := false
		OnFalse(&val, func() {
			executed = true
		})
		if !executed {
			t.Error("Should have executed")
		}
	})

	t.Run("does not execute with true", func(t *testing.T) {
		val := true
		OnFalse(&val, func() {
			t.Error("Should not execute")
		})
	})

	t.Run("safe with nil pointer", func(t *testing.T) {
		OnFalse(nil, func() {
			t.Error("Should not execute")
		})
	})
}

func TestExecIfSetHelper(t *testing.T) {
	t.Run("executes with valid function", func(t *testing.T) {
		executed := false
		fn := func() {
			executed = true
		}
		ExecIfSet(&fn)
		if !executed {
			t.Error("Should have executed")
		}
	})

	t.Run("safe with nil pointer", func(t *testing.T) {
		ExecIfSet(nil)
	})

	t.Run("safe with nil function", func(t *testing.T) {
		var fn func()
		ExecIfSet(&fn)
	})
}
