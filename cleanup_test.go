package cleanup

import (
	"errors"
	"strings"
	"testing"
)

const (
	failedStr  = "failed"
	successStr = "success"
	alwaysStr  = "always"
)

var (
	errFailed = errors.New("I have failed")
)

// ExampleFunction just an example function
func ExampleFunction(fail bool) (usefulthing *string, err error) {
	cleaner := NewCleaner(&err)
	defer cleaner.Clean()

	usefulthing = new(string)

	// (some logic that needs to be executed only if MyFunc returns an error)
	cleaner.OnError(func(e error) {
		if e == errFailed {
			*usefulthing += ":" + failedStr
		} else {
			panic("not my err")
		}
	})
	// practical context use case:
	//
	// 	cleaner.AddOnError(func(e err) {
	// 		if e == ErrTxFailed {
	// 			transaction.Rollback()
	// 		} else {
	// 			log.Errorf("failed request ID: %s", reqID)
	// 		}
	// 	})

	// (some logic that needs to be executed only if ExampleFunction returns no error)
	cleaner.OnNil(func() {
		*usefulthing += ":" + successStr
		uerr := errors.New("updated error")
		*cleaner.Errptr = uerr
	})
	// practical context use case:
	//
	// 	cleaner.AddOnNil(func() { log.Info("everything went ok!") })

	// (some logic that needs to be executed always, no matter what)
	cleaner.Always(func() {
		*usefulthing += ":" + alwaysStr
	})
	// practical context use case:
	//
	// 	cleaner.AddAlways(func() { wipebytes(thisByteSliceShouldBeZeroedAlways) })

	if fail {
		err = errFailed
	}

	return
}

func Test(t *testing.T) {
	thing, err := ExampleFunction(true)

	if err == nil {
		t.Error("err should return an error")
	}

	things := strings.Split(*thing, ":")

	if things[1] != alwaysStr {
		t.Errorf("first thing should be '%s' instead is '%s'", alwaysStr, things[1])
	}

	if things[2] != failedStr {
		t.Errorf("second thing should be '%s' instead is '%s'", failedStr, things[2])
	}

}
