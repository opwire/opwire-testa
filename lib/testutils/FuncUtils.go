package testutils

func GetError(out interface{}, err error) error {
	return err
}

func GetFirstResult_bool(out bool, err error) bool {
	return out
}

func GetFirstResult_int(out int, err error) int {
	return out
}

func GetFirstResult_string(out string, err error) string {
	return out
}

func GetFirstResult_any(out interface{}, err error) interface{} {
	return out
}
