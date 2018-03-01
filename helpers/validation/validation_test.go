package validation

import "testing"

func TestValidateEmail(t *testing.T) {
	if ValidateEmail("") {
		t.Error(`ValidateEmail("") == true`)
	}

	if ValidateEmail("ninhnt") {
		t.Error(`ValidateEmail("ninhnt") == true`)
	}

	if ValidateEmail(" @.") {
		t.Error(`ValidateEmail(" @.") == true`)
	}

	if !ValidateEmail("ninhnt@thebigdev.com") {
		t.Error(`ValidateEmail("ninhnt@thebigdev.com") == false`)
	}
}
