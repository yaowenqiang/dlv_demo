package main

import (
	"testing"
)
//run test|debug test
func TestPasswordIsCorrectLength(t *testing.T) {
	const length = 42
	pwd := genPassword(length, false)
	if len(pwd) != length + 1 {
		t.Fatal("Password is not correct length")
	}
}
