# cleanup

[![Godoc Reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/alvarolm/cleanup)
[![Go Report](https://img.shields.io/badge/go%20report-A%2B-brightgreen.svg)](https://goreportcard.com/report/github.com/alvarolm/cleanup)

[![donate](https://img.shields.io/badge/donate-a%20bus%20ticket%2C%20cup%20of%20coffe%2C%20anything%20you%20can%2C%20thanks!-orange.svg)](https://www.paypal.com/cgi-bin/webscr?cmd=_xclick&business=alvarofleivam%40gmail%2ecom&lc=AL&item_name=Donation%20%5b%20for%20a%20bus%20ticket%2c%20coffe%20anything%20you%20can%20I%27m%20happy%20thanks%20%21%20%3a%29%20%5d&item_number=donation&button_subtype=services&currency_code=USD&bn=PP%2dBuyNowBF%3abtn_buynowCC_LG%2egif%3aNonHosted)

deferred function utilities
cleans up for you !

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
