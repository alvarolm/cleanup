# cleanup

[![Godoc Reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/alvarolm/cleanup)
[![Go Report](https://img.shields.io/badge/go%20report-A%2B-brightgreen.svg)](https://goreportcard.com/report/github.com/alvarolm/cleanup)

deferred function utilities

example:
```go
func ProcessQuery()  (err error) {
	cleaner := NewCleaner(&err)
	defer cleaner.Clean()

	...

	// (some logic that needs to be executed only if ProcessQuery returns an error)
	cleaner.OnError(func(e err) {
		if e == ErrTxFailed {
			transaction.Rollback()
		} else {
			log.Errorf("failed request ID: %s", reqID)
		}
	})

	... 

	// (some logic that needs to be executed only if ProcessQuery returns no error)
	cleaner.OnNil(func() { log.Infof("request ID: %s succeded", reqID)

	...

	// (some logic that needs to be executed always no matter what)
	cleaner.Always(func() { wipebytes(thisByteSliceShouldBeZeroedAlways) })

	...

	return
}
```
