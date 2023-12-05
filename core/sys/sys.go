package sys

import "os"

func Exit(exitCode ...int) {
	wg.Wait()
	var code int
	if len(exitCode) > 0 {
		code = exitCode[0]
	}
	os.Exit(code)
}

func Abort(err error) {
	cancel(err)
}
