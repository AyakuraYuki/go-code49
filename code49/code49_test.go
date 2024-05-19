package code49

import (
	"testing"
)

func TestDecode(t *testing.T) {
	codes := []string{
		"11143121314115211131114321124131314",
		"11221611211411251111225122311314214",
		"11123232212411212332131231332321114",
		"11251311211242114112215212413213114",
		"11123121511212521211113243422213114",
		"11224211311211313421211153141112154",
	}
	t.Log(Decode(codes, true))
	t.Log(Decode(codes, false))
}

func TestDecodeRaw(t *testing.T) {
	codes := []string{
		"11143121314115211131114321124131314",
		"11221611211411251111225122311314214",
		"11123232212411212332131231332321114",
		"11251311211242114112215212413213114",
		"11123121511212521211113243422213114",
		"11224211311211313421211153141112154",
	}
	t.Log(DecodeRaw(codes))
}
