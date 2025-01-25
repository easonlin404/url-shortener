package utils

import "testing"

func TestBase62Encode_Zero(t *testing.T) {
	result := Base62Encode(0)
	expected := "0"
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestBase62Encode_SingleDigit(t *testing.T) {
	result := Base62Encode(1)
	expected := "1"
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestBase62Encode_MultipleDigits(t *testing.T) {
	result := Base62Encode(12345)
	expected := "3d7"
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestBase62Encode_LargeNumber(t *testing.T) {
	result := Base62Encode(9876543210)
	expected := "aMoY42"
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestBase62Encode_NegativeNumber(t *testing.T) {
	result := Base62Encode(-12345)
	expected := ""
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestBase62Encode_ExactBase62(t *testing.T) {
	result := Base62Encode(62)
	expected := "10"
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}
