package entity_test

import (
	"testing"

	"github.com/AnatoliyBr/todo-app/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestList_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		title   string
		isValid bool
	}{
		{
			name:    "valid",
			title:   "JOB",
			isValid: true,
		},
		{
			name:    "mixedcase with whitespace",
			title:   "development    OF THE Zhivoi Zvuk Club 2 ",
			isValid: true,
		},
		{
			name:    "invalid symbols",
			title:   "ITMO ?#@*&%!",
			isValid: false,
		},
		{
			name:    "empty",
			title:   "",
			isValid: false,
		},
		{
			name:    "long title",
			title:   "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := &entity.List{ListTitle: tc.title}
			if tc.isValid {
				assert.NoError(t, l.Validate())
			} else {
				assert.Error(t, l.Validate())
			}
		})
	}
}
