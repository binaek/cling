package cling

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

type configTarget struct {
	valType   reflect.Type
	structIdx int
}

type configTargets map[string]configTarget

// Hydrate populates the destination struct based on command-line arguments and context.
func Hydrate[T any](ctx context.Context, argArguments []string, destination *T) error {
	if destination == nil {
		return errors.New("destination cannot be nil")
	}
	cmd, ok := CommandFromContext(ctx)
	if !ok {
		return errors.New("invalid state - context is not derived from CLIng supplied context")
	}

	// parse the arguments
	argFlags, argArguments := parseArguments(argArguments)
	targets, err := extractConfigTargets(destination)
	if err != nil {
		return errors.Wrap(err, "failed to extract config targets")
	}

	// make sure that we have targets for all required
	for _, cmdFlag := range cmd.flags {
		if !cmdFlag.isRequired() {
			continue
		}
		if _, found := targets[cmdFlag.Name()]; !found {
			return errors.Errorf("could not find target for required flag: '%s'", cmdFlag.Name())
		}
	}

	for _, cmdArg := range cmd.arguments {
		if !cmdArg.isRequired() {
			continue
		}
		if _, found := targets[cmdArg.Name()]; !found {
			return errors.Errorf("could not find target for required argument: '%s'", cmdArg.Name())
		}
	}

	destVal := reflect.ValueOf(destination).Elem()

	if err := hydrateFlags(cmd, argFlags, destVal, targets); err != nil {
		return err
	}

	if err := hydrateArgs(cmd, argArguments, destVal, targets); err != nil {
		return err
	}

	return nil
}

func hydrateArgs(cmd *Command, args []string, destination reflect.Value, targets configTargets) error {
	// verify we have at least the required number of arguments
	requiredArguments := 0
	for _, argument := range cmd.arguments {
		if argument.isRequired() {
			requiredArguments++
		}
	}

	if len(args) < requiredArguments {
		return errors.Errorf("missing at least one required argument. need '%d' - got '%d'", requiredArguments, len(args))
	}

	for idx, argument := range cmd.arguments {
		validator := NoOpValidator()
		if flagWithValidator, ok := argument.(CmdInputWithDefaultAndValidator[any]); ok {
			validator = flagWithValidator.getValidator().(Validator[any])
		}

		target, ok := targets[argument.Name()]
		field := destination.Field(target.structIdx)
		if idx < len(args) {
			if !ok {
				return errors.Errorf("could not find target for '%s'", argument.Name())
			}
			if err := setFieldFromString(field, args[idx], validator); err != nil {
				return errors.Wrapf(err, "failed to set argument '%s'", argument.Name())
			}
		} else {
			// go with default
			val := fmt.Sprint(argument.getDefault())
			if ok {
				// put in the default
				if err := setFieldFromString(field, val, validator); err != nil {
					return errors.Wrapf(err, "failed to set argument '%s'", argument.Name())
				}
			}
		}
	}

	return nil
}

func hydrateFlags(cmd *Command, v map[string][]string, destination reflect.Value, targets configTargets) error {
	// get defined flags
	for _, flag := range cmd.flags {
		name := flag.Name()
		values, ok := v[name]

		validator := NoOpValidator()
		if flagWithValidator, ok := flag.(ValidatorProvider); ok && flagWithValidator.getValidator() != nil {
			validator = flagWithValidator.getValidator()
		}

		if !ok && flag.hasDefault() {
			// get the default
			def := flag.getDefault()
			// run it through the validator
			if err := validator.Validate(def); err != nil {
				return errors.Wrapf(err, "cannot set invalid default '%v' for '%s'", def, name)
			}
			// if def is a slice
			if reflect.TypeOf(def).Kind() == reflect.Slice {
				values = []string{}
				for i := 0; i < reflect.ValueOf(def).Len(); i++ {
					values = append(values, fmt.Sprint(reflect.ValueOf(def).Index(i).Interface()))
				}
			} else {
				values = []string{fmt.Sprint(def)}
			}
		}

		if len(values) == 0 {
			// try to populate from env
			for _, envKey := range flag.envSources() {
				if val, ok := os.LookupEnv(envKey); ok {
					values = []string{val}
				}
			}
		}

		if (flag.isRequired()) && (!ok || len(values) == 0) {
			return errors.Errorf("missing required flag '%s'", flag.Name())
		}

		target, ok := targets[name]
		if !ok {
			return errors.Errorf("could not find target for '%s'", name)
		}
		field := destination.Field(target.structIdx)

		if !field.IsValid() {
			return errors.Errorf("no valid field found for flag '%s'", name)
		}
		if !field.CanSet() {
			return errors.Errorf("field for flag '%s' cannot be set", name)
		}
		if field.Kind() == reflect.Slice {
			for _, valueStr := range values {
				if err := setFieldFromString(field, valueStr, validator); err != nil {
					return errors.Wrapf(err, "failed to set flag '%s'", name)
				}
			}
		} else {
			valueStr := values[0]
			if err := setFieldFromString(field, valueStr, validator); err != nil {
				return errors.Wrapf(err, "failed to set flag '%s'", name)
			}
		}
	}

	return nil
}

func parseArguments(args []string) (flags map[string][]string, arguments []string) {
	flags = make(map[string][]string)
	arguments = make([]string, 0)
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
			arguments = append(arguments, arg)
		}
	}

	return flags, arguments
}

func extractConfigTargets(config any) (targets map[string]configTarget, e error) {
	targets = make(map[string]configTarget)
	configType := reflect.TypeOf(config)
	if configType.Kind() != reflect.Ptr || configType.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("CLIng can only parse command line arguments into structs, got %v", configType.Kind())
	}

	v := reflect.ValueOf(config).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		nameTag, ok := field.Tag.Lookup("cling-name")
		if !ok {
			// this is not a field we are interested in
			continue
		}
		target := configTarget{
			valType:   field.Type,
			structIdx: i,
		}
		if _, ok := targets[nameTag]; ok {
			return nil, errors.Errorf("found duplicate 'cling:name' in config")
		}
		targets[nameTag] = target
	}

	return targets, nil
}
