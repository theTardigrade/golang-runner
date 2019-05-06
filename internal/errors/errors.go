package errors

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func Judge(err error) string {
	if err == nil {
		return "SUCCESS"
	}

	return "FAILURE"
}
