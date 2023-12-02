package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	NestedTable struct {
		Name                    string     `validate:"in:test"`
		ResponseLine            Response   `validate:"nested"`
		ResponseArray           []Response `validate:"nested"`
		ResponseArrayWOValidate []Response
	}

	BadTagValue struct {
		UndefinedTags  string `validate:"udTest:13|SomeNew:123,456"`
		BadStringsTags string `validate:"len:bad|in|regexp:\\bad"`
		BadIntTags     int    `validate:"min:bad|max:bad|in:test,some"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name           string
		in             interface{}
		expectedSErr   error
		expectedVError error
	}{
		{
			name:         "nil input",
			in:           nil,
			expectedSErr: errSystem,
		},
		{
			name:         "non struct input",
			in:           42,
			expectedSErr: errSystem,
		},
		{
			name: "struct without tags",
			in: Token{
				Header:    nil,
				Payload:   nil,
				Signature: nil,
			},
			expectedSErr:   nil,
			expectedVError: nil,
		},
		{
			name: "good system error",
			in: BadTagValue{
				UndefinedTags:  "test",
				BadStringsTags: "test",
				BadIntTags:     1234,
			},
			expectedSErr: errSystem,
		},
		{
			name: "good check user",
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "any name",
				Age:    18,
				Email:  "Test@test.com",
				Role:   "admin",
				Phones: []string{"89181113322"},
				meta:   []byte("any date"),
			},
			expectedSErr:   nil,
			expectedVError: nil,
		},
		{
			name: "good check nested",
			in: NestedTable{
				Name: "test",
				ResponseLine: Response{
					Code: 200,
				},
				ResponseArray: []Response{
					{Code: 404},
					{Code: 500},
				},
				ResponseArrayWOValidate: []Response{
					{Code: -1},
				},
			},
			expectedSErr:   nil,
			expectedVError: nil,
		},
		{
			name: "good validation user",
			in: User{
				ID:     "1234",
				Name:   "any name",
				Age:    15,
				Email:  "wrong",
				Role:   "-",
				Phones: []string{"000"},
				meta:   []byte("any date"),
			},
			expectedSErr:   nil,
			expectedVError: errValidation,
		},
		{
			name: "good validation app",
			in: App{
				Version: "123467",
			},
			expectedSErr:   nil,
			expectedVError: errValidation,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d: %v", i, tt.name), func(t *testing.T) {
			tt := tt
			err := Validate(tt.in)
			fmt.Printf("%v\n", err)
			if tt.expectedSErr != nil {
				if err == nil {
					t.Errorf("expected system error but expected nil")
					return
				}
				if !(errors.Is(err, tt.expectedSErr)) {
					t.Errorf("error '%v' but expected '%v'", err, tt.expectedSErr)
				}
				return
			}
			if tt.expectedVError != nil {
				if err == nil {
					t.Errorf("expected validation errors but expected nil")
					return
				}
				if !errors.Is(err, tt.expectedVError) {
					t.Errorf("expected validation error %v but got %v", tt.expectedVError, err)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error '%v'", err)
			}
			_ = tt
		})
	}
}
