package validator

import (
	"testing"
)

func TestValidateEmail(t *testing.T) {
	t.Run("valid 1", func(t_ *testing.T) {
		if err := ValidateEmail("glob.den@gmail.com"); err != nil {
			t_.Errorf("Error: %s", err.Error())
		} else {
			t_.Logf("Success - no errors")
		}
	})

	t.Run("valid 2", func(t_ *testing.T) {
		if err := ValidateEmail("abrakadabra@yandex.ru"); err != nil {
			t_.Errorf("Error: %s", err.Error())
		} else {
			t_.Logf("Success - no errors")
		}
	})

	t.Run("valid 3", func(t_ *testing.T) {
		if err := ValidateEmail("fatCat89@yandex.ru"); err != nil {
			t_.Errorf("Error: %s", err.Error())
		} else {
			t_.Logf("Success - no errors")
		}
	})

	t.Run("invalid - double @", func(t_ *testing.T) {
		if err := ValidateEmail("glob@den@gmail.com"); err != nil {
			t_.Logf("Success - error found as it expected - %s", err.Error())
		} else {
			t_.Errorf("Error: no errors, but it should be")
		}
	})

	t.Run("invalid - no @", func(t_ *testing.T) {
		if err := ValidateEmail("abrakadabrayandex.ru"); err != nil {
			t_.Logf("Success - error found as it expected - %s", err.Error())
		} else {
			t_.Errorf("Error: no errors, but it should be")
		}
	})

	t.Run("invalid - @ as last symbol", func(t_ *testing.T) {
		if err := ValidateEmail("abrakadabra@"); err != nil {
			t_.Logf("Success - error found as it expected - %s", err.Error())
		} else {
			t_.Errorf("Error: no errors, but it should be")
		}
	})

	t.Run("invalid - only dots", func(t_ *testing.T) {
		if err := ValidateEmail("....@...."); err != nil {
			t_.Logf("Success - error found as it expected - %s", err.Error())
		} else {
			t_.Errorf("Error: no errors, but it should be")
		}
	})

	t.Run("invalid - length", func(t_ *testing.T) {
		if err := ValidateEmail("a"); err != nil {
			t_.Logf("Success - error found as it expected - %s", err.Error())
		} else {
			t_.Errorf("Error: no errors, but it should be")
		}
	})
}
