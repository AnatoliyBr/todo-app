package entity

import (
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
)

type List struct {
	ListID    int    `json:"list_id"`
	ListTitle string `json:"list_title"`
	UserID    int    `json:"user_id,omitempty"`
}

func (l *List) Validate() error {
	l.ListTitle = strings.Join(strings.Fields(l.ListTitle), " ")
	l.ListTitle = strings.ToUpper(l.ListTitle)

	return validation.ValidateStruct(
		l,
		validation.Field(
			&l.ListTitle,
			validation.Required,
			validation.Match(regexp.MustCompile(`^[0-9A-Z ]+$`)),
			validation.Length(0, 50)),
	)
}
