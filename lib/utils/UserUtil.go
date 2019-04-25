package utils

import (
	"os/user"
)

func FindUsername() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.Username, nil
}
