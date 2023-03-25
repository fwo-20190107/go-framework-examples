package validation

import (
	"errors"
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var NgWord = &ngWord{message: "contains prohibited words"}

type ngWord struct {
	message string
}

func (v *ngWord) Validate(value interface{}) error {
	isString, str, _, _ := validation.StringOrBytes(value)
	if !isString {
		return errors.New("value is not a comparable string")
	}

	for _, word := range ngWords {
		if matched, err := regexp.MatchString(word, str); err != nil {
			return err
		} else if matched {
			return fmt.Errorf(v.message+":%s", word)
		}
	}
	return nil
}

func (v *ngWord) Error(message string) *ngWord {
	return &ngWord{
		message: message,
	}
}

var ngWords = []string{
	"hoge",
	"fuga",
}

var _ validation.Rule = (*ngWord)(nil)
