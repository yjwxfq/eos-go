package eos

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/streamingfast/validator"
	"github.com/thedevsaddam/govalidator"
)

func init() {
	govalidator.AddCustomRule("eos.blockNum", EOSBlockNumRule)
	govalidator.AddCustomRule("eos.name", EOSNameRule)
}

func EOSBlockNumRule(field string, rule string, message string, value interface{}) error {
	val, ok := value.(string)
	if !ok {
		return fmt.Errorf("The %s field must be a string", field)
	}

	_, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return fmt.Errorf("The %s field must be a valid EOS block num", field)
	}

	return nil
}

func EOSNameRule(field string, rule string, message string, value interface{}) error {
	checkName := func(field string, name string) error {
		if !IsValidName(name) {
			return fmt.Errorf("The %s field must be a valid EOS name", field)
		}

		return nil
	}

	switch v := value.(type) {
	case string:
		return checkName(field, v)
	case Name, PermissionName, ActionName, AccountName, TableName:
		return checkName(field, fmt.Sprintf("%s", v))
	default:
		return fmt.Errorf("The %s field is not a known type for an EOS name", field)
	}
}

func EOSExtendedNameRule(field string, rule string, message string, value interface{}) error {
	checkName := func(field string, name string) error {
		if !IsValidExtendedName(name) {
			return fmt.Errorf("The %s field must be a valid EOS name", field)
		}

		return nil
	}

	switch v := value.(type) {
	case string:
		return checkName(field, v)
	case Symbol:
		return checkName(field, v.String())
	case SymbolCode:
		return checkName(field, v.String())
	case Name, PermissionName, ActionName, AccountName, TableName:
		return checkName(field, fmt.Sprintf("%s", v))
	default:
		return fmt.Errorf("The %s field is not a known type for an EOS name", field)
	}
}

func EOSNamesListRuleFactory(sep string, maxCount int) validator.Rule {
	return validator.StringListRuleFactory(sep, maxCount, EOSNameRule)
}

func EOSExtendedNamesListRuleFactory(sep string, maxCount int) validator.Rule {
	return validator.StringListRuleFactory(sep, maxCount, EOSExtendedNameRule)
}

func EOSTrxIDRule(field string, rule string, message string, value interface{}) error {
	err := validator.HexRule(field, rule, message, value)
	if err != nil {
		return err
	}

	val := value.(string)
	if len(val) != 64 {
		return fmt.Errorf("The %s field must have exactly 64 characters", field)
	}

	return nil
}

var symbolRegexp = regexp.MustCompile(`^[0-9],[A-Z]{1,7}$`)
var symbolCodeRegexp = regexp.MustCompile(`^[A-Z]{1,7}$`)
var nameRegexp = regexp.MustCompile(`^[\.a-z1-5]{0,13}$`)

func ExplodeNames(input string, sep string) (names []string) {
	rawNames := strings.Split(input, sep)
	for _, rawName := range rawNames {
		account := strings.TrimSpace(rawName)
		if account == "" {
			continue
		}

		names = append(names, rawName)
	}

	return
}

func IsValidName(input string) bool {
	// An empty string name means a uint64 transformed name with a 0 value
	if input == "" {
		return true
	}

	return nameRegexp.MatchString(input)
}

func IsValidExtendedName(input string) bool {
	// An empty string name means a uint64 transformed name with a 0 value
	if input == "" {
		return true
	}

	return nameRegexp.MatchString(input) || symbolCodeRegexp.MatchString(input) || symbolRegexp.MatchString(input)
}
