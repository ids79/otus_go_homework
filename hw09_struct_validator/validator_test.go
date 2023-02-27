package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Code   int      `validate:"in:200,404,500"`
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	Ex1 struct {
		Version string `validate:"min:18"`
	}

	Ex2 struct {
		Version string `validate:"max:18"`
	}

	Ex3 struct {
		Art int `validate:"regexp:^\\w+$"`
	}

	Ex4 struct {
		Num int `validate:"len:18"`
	}

	Token struct {
		ID     string   `json:"id" validate:"len:36"`
		Header Response `validate:"nested"`
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

var tests = []struct {
	in          interface{}
	field       string
	expectedErr error
}{
	{
		in: User{
			ID:     "12234567876543212345",
			Name:   "name",
			Code:   200,
			Age:    26,
			Email:  "1@ya.ru",
			Role:   UserRole("admin"),
			Phones: []string{"12312131523"},
			meta:   []byte("1dfdfsd"),
		},
		field:       "",
		expectedErr: nil,
	},
	{
		in: User{
			ID:     "12234567876543212345",
			Name:   "name",
			Code:   200,
			Age:    16,
			Email:  "1@ya.ru",
			Role:   UserRole("admin"),
			Phones: []string{"12312131523"},
			meta:   []byte("1dfdfsd"),
		},
		field:       "Age",
		expectedErr: ErrValueLessMinimum,
	},
	{
		in: User{
			ID:     "12234567876543212345",
			Name:   "name",
			Code:   200,
			Age:    80,
			Email:  "1@ya.ru",
			Role:   UserRole("admin"),
			Phones: []string{"12312131523"},
			meta:   []byte("1dfdfsd"),
		},
		field:       "Age",
		expectedErr: ErrValueMoreMaximum,
	},
	{
		in: User{
			ID:     "1223456787654321234545678765434567898",
			Name:   "name",
			Code:   200,
			Age:    26,
			Email:  "1@ya.ru",
			Role:   UserRole("admin"),
			Phones: []string{"12312131523"},
			meta:   []byte("1dfdfsd"),
		},
		field:       "ID",
		expectedErr: ErrStringLengthLongerAllowed,
	},
	{
		in: User{
			ID:     "12234567876543212345",
			Name:   "name",
			Code:   200,
			Age:    26,
			Email:  "1@yaru",
			Role:   UserRole("admin"),
			Phones: []string{"12312131523"},
			meta:   []byte("1dfdfsd"),
		},
		field:       "Email",
		expectedErr: ErrValueNotMatchPattern,
	},
	{
		in: User{
			ID:     "12234567876543212345",
			Name:   "name",
			Code:   200,
			Age:    26,
			Email:  "1@ya.ru",
			Role:   UserRole("manager"),
			Phones: []string{"12312131523"},
			meta:   []byte("1dfdfsd"),
		},
		field:       "Role",
		expectedErr: ErrValueNotIncludedAllowedList,
	},
	{
		in: User{
			ID:     "12234567876543212345",
			Name:   "name",
			Code:   100,
			Age:    26,
			Email:  "1@ya.ru",
			Role:   UserRole("admin"),
			Phones: []string{"12312131523"},
			meta:   []byte("1dfdfsd"),
		},
		field:       "Code",
		expectedErr: ErrValueNotIncludedAllowedList,
	},
	{
		in: User{
			ID:     "12234567876543212345",
			Name:   "name",
			Code:   200,
			Age:    26,
			Email:  "1@ya.ru",
			Role:   UserRole("admin"),
			Phones: []string{"12312131523", "123456787654"},
			meta:   []byte("1dfdfsd"),
		},
		field:       "Phones",
		expectedErr: ErrStringLengthLongerAllowed,
	},
	{
		in: Token{
			ID: "12213123",
			Header: Response{
				Code: 100,
				Body: "",
			},
		},
		field:       "Code",
		expectedErr: ErrValueNotIncludedAllowedList,
	},
	{
		in:          Ex1{"dddddd"},
		field:       "Version",
		expectedErr: nil,
	},
	{
		in:          Ex2{"dddddd"},
		field:       "Version",
		expectedErr: nil,
	},
	{
		in:          Ex3{122133123},
		field:       "Art",
		expectedErr: nil,
	},
	{
		in:          Ex4{122133123},
		field:       "Num",
		expectedErr: nil,
	},
	{
		in:          10,
		field:       "",
		expectedErr: nil,
	},
}

func TestValidate(t *testing.T) {
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d field %s", i, tt.field), func(t *testing.T) {
			tt := tt
			t.Parallel()
			var rez interface{} = Validate(tt.in)
			switch err := rez.(type) {
			case ValidationErrors:
				if len(err) == 0 {
					require.Nil(t, tt.expectedErr)
					require.Equal(t, tt.field, "")
				} else {
					for _, valid := range err {
						require.Equal(t, tt.field, valid.Field)
						require.ErrorIs(t, tt.expectedErr, valid.Err)
					}
				}
			case error:
				require.NotNil(t, err)
				fmt.Println(err.Error())
			default:
			}
		})
	}
}
