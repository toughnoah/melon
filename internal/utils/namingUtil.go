package utils

import (
	"errors"
	"fmt"
	"regexp"
)

const DefaultMatchAllExpr = ".*?"

type Config struct {
	Deployment *Deployment `yaml:"deployment"`
}

type Deployment struct {
	Image  *string `yaml:"image"`
	Naming *string `yaml:"naming"`
}

func validateNaming(name, expr string) error {
	reg, err := regexp.Compile(expr)
	if err != nil {
		return err
	}

	match := reg.FindAllStringSubmatch(name, -1)
	if match != nil {
		return nil
	}
	return errors.New(fmt.Sprintf(matchExprError, name, expr))
}

func ValidateNaming(name, path, kind string) error {
	expr, err := GetValFromConfig(path, kind)
	if err != nil {
		return err
	}
	exprStr, ok := expr.(string)
	if !ok {
		return errors.New(fmt.Sprintf(badValueTypeError, kind, "string"))
	}
	if err = validateNaming(name, exprStr); err != nil {
		return err
	}
	return nil
}
