package cling

import (
	"fmt"
	"reflect"

	"github.com/gertd/go-pluralize"
)

type Arg struct {
	name        string
	description string
	position    int
	argType     reflect.Type
	structIdx   int
	optional    bool
}

type Flag struct {
	name        string
	description string
	shorthand   rune
	flagType    reflect.Type
	structIdx   int
	optional    bool
}

func (f *Flag) getHelpType() string {
	switch f.flagType.Kind() {
	case reflect.Slice:
		var elem string
		switch f.flagType.Elem().Kind() {
		case reflect.Bool:
			elem = "booleans [true|false]"
		case reflect.Pointer:
			elem = f.flagType.Elem().Kind().String()
		default:
			elem = f.flagType.Elem().Kind().String()
		}
		pluralize := pluralize.NewClient()
		elem = pluralize.Plural(elem)
		elem = fmt.Sprintf("%s (comma separated)", elem)
		return elem
	case reflect.Bool:
		return "boolean [true|false]"
	case reflect.Pointer:
		return f.flagType.Elem().Kind().String()
	default:
		return f.flagType.String()
	}
}
