package cling

import "testing"

func TestHydrate(t *testing.T) {

}

func TestParseArguments(t *testing.T) {
	args := []string{"arg1", "arg2", "--flag1", "value1", "--flag2", "value2"}
	flags, args := parseArguments(args)
	if len(flags) != 2 {
		t.Errorf("Expected 2 flags, got %d", len(flags))
	}
	if len(args) != 2 {
		t.Errorf("Expected 2 args, got %d", len(args))
	}
}
