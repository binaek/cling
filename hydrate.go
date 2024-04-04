package cling

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func Hydrate[T any](ctx context.Context, args []string, destination *T) error {
	// parse the arguments
	flags, args := parseArguments(args)
	defFlags, defArgs, err := extractConfigDefinitions(destination)
	if err != nil {
		return err
	}

	// verify that there are no non-optional arguments after an optional one
	for idx, arg := range defArgs {
		if arg.optional && idx < len(defArgs)-1 && !defArgs[idx+1].optional {
			return fmt.Errorf("non-optional argument %s after optional argument %s", defArgs[idx+1].description, arg.description)
		}
	}

	destVal := reflect.ValueOf(destination).Elem()

	for _, flag := range defFlags {
		flagValues, ok := flags[flag.name]
		if !ok || len(flagValues) == 0 {
			if !flag.optional {
				return fmt.Errorf("missing required flag %s", flag.name)
			}
			continue // No values provided for this flag
		}
		field := destVal.Field(flag.structIdx)
		if !field.IsValid() || !field.CanSet() {
			return fmt.Errorf("no field found for flag '%s'", flag.name)
		}

		if field.Kind() == reflect.Slice {
			for _, valueStr := range flagValues {
				if err := setFieldFromString(field, valueStr); err != nil {
					return errors.Wrapf(err, "failed to set flag '%s'", flag.name)
				}
			}
		} else {
			valueStr := flagValues[0]
			if err := setFieldFromString(field, valueStr); err != nil {
				return errors.Wrapf(err, "failed to set flag '%s'", flag.name)
			}
		}
	}

	// Populate positionals
	for i, p := range defArgs {
		if i >= len(args) {
			continue // Not enough positional arguments provided
		}

		field := destVal.Field(p.structIdx)
		if !field.IsValid() || !field.CanSet() {
			return fmt.Errorf("no field found for positional argument '%s'", p.name)
		}

		valueStr := args[i]
		if err := setFieldFromString(field, valueStr); err != nil {
			return errors.Wrapf(err, "failed to set positional argument '%s'", p.name)
		}
	}

	return nil
}

func parseArguments(args []string) (map[string][]string, []string) {
	flags := make(map[string][]string)
	var positionals []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--") {
			parts := strings.SplitN(arg, "=", 2)
			flagName := strings.TrimPrefix(parts[0], "--")
			if len(parts) == 2 {
				// Handle --flag=value
				flags[flagName] = append(flags[flagName], parts[1])
			} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
				// Handle --flag value
				flags[flagName] = append(flags[flagName], args[i+1])
				i++ // Skip the next element as it is already used as a value
			} else {
				// Handle flags without values, assign an empty string or a default value
				flags[flagName] = append(flags[flagName], "")
			}
		} else {
			positionals = append(positionals, arg)
		}
	}

	return flags, positionals
}

func extractConfigDefinitions(config any) ([]*Flag, []*Arg, error) {
	var flags []*Flag
	var positionals []*Arg

	configType := reflect.TypeOf(config)
	if configType.Kind() != reflect.Ptr || configType.Elem().Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("CLIng can only parse command line arguments into structs, got %v", configType.Kind())
	}

	v := reflect.ValueOf(config).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldType := t.Field(i)
		if arg, err := getPositionalFromStructField(fieldType, i); err == nil && arg != nil {
			positionals = append(positionals, arg)
		} else if flag, err := getFlagFromStructField(fieldType, i); err == nil && flag != nil {
			flags = append(flags, flag)
		} else if err != nil {
			return nil, nil, err
		}
	}

	slices.SortStableFunc(positionals, func(left *Arg, right *Arg) int {
		return left.position - right.position
	})

	return flags, positionals, nil
}

func getFlagFromStructField(field reflect.StructField, index int) (*Flag, error) {
	tag, ok := field.Tag.Lookup("name")
	if !ok {
		return nil, fmt.Errorf("flag field must have a name tag")
	}

	flag := &Flag{
		name:        tag,
		description: "",
		shorthand:   0,
		flagType:    field.Type,
		structIdx:   index,
	}

	description, ok := field.Tag.Lookup("description")
	if ok {
		flag.description = description
	}

	short, ok := field.Tag.Lookup("short")
	if len(short) != 1 {
		return nil, fmt.Errorf("short flag must be a single character")
	}
	if ok {
		flag.shorthand = rune(short[0])
	}

	if field.Type.Kind() == reflect.Pointer {
		flag.optional = true
	}

	return flag, nil
}

func getPositionalFromStructField(field reflect.StructField, index int) (*Arg, error) {
	tag, ok := field.Tag.Lookup("position")
	if !ok {
		return nil, nil
	}

	position, err := strconv.Atoi(tag)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert positional tag to int: '%s'", tag)
	}

	arg := &Arg{
		description: "",
		position:    position,
		argType:     field.Type,
		structIdx:   index,
		name:        field.Tag.Get("name"),
	}

	description, ok := field.Tag.Lookup("description")
	if ok {
		arg.description = description
	}

	if field.Type.Kind() == reflect.Pointer {
		arg.optional = true
	}

	return arg, nil
}
