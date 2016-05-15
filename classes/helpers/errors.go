package helpers

// CheckErrAndPanic speaks for itself
func CheckErrAndPanic(err error) {
	if err != nil {
		panic(err.Error())
	}
}
