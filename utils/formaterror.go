package utils

import "errors"

func FormatError(t string) error {
	switch {
	case t == "supabase_signUp":
		return errors.New("email is already existed")
	default:
		return errors.New("invalid argument")
	}
}
