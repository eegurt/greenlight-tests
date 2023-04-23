package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"greenlight.bcc/internal/assert"
)

func TestCreateUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	const (
		validUsername  = "Bob The Tester"
		validEmail     = "bob_tester@example.com"
		duplicateEmail = "bob@example.com"
		validPassword  = "Secret-Password"
	)

	initData := struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Name:     "Bob The First",
		Email:    duplicateEmail,
		Password: validPassword,
	}

	initBytes, err := json.Marshal(&initData)
	if err != nil {
		t.Fatalf("Failed to init data model: %v", err)
	}

	ts.postForm(t, "/v1/users", initBytes)

	tests := []struct {
		name     string
		Username string
		Email    string
		Password string
		wantCode int
	}{
		{
			name:     "Valid user registration",
			Username: validUsername,
			Email:    validEmail,
			Password: validPassword,
			wantCode: http.StatusCreated,
		},
		{
			name:     "Wrong input",
			Username: validUsername,
			Email:    validEmail,
			Password: validPassword,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Empty username",
			Username: "",
			Email:    validEmail,
			Password: validPassword,
			wantCode: http.StatusUnprocessableEntity,
		},
		// {
		// 	name:     "Creating user with existing email",
		// 	Username: "Bob The Second",
		// 	Email:    duplicateEmail,
		// 	Password: validPassword,
		// 	wantCode: http.StatusUnprocessableEntity,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputData := struct {
				Name     string `json:"name"`
				Email    string `json:"email"`
				Password string `json:"password"`
			}{
				Name:     tt.Username,
				Email:    tt.Email,
				Password: tt.Password,
			}

			b, err := json.Marshal(&inputData)
			if err != nil {
				t.Fatal("wrong input data")
			}
			if tt.name == "Wrong input" {
				b = append(b, 'a')
			}

			code, _, _ := ts.postForm(t, "/v1/users", b)

			assert.Equal(t, code, tt.wantCode)
		})
	}
}
