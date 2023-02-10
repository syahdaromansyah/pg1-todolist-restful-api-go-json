package helper

func DoPanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
