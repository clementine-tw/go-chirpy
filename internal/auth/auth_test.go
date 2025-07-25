package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {

	password1 := "currentPass123"
	passwordHash1, _ := HashPassword(password1)

	password2 := "anotherPass456"
	passwordHash2, _ := HashPassword(password2)

	cases := []struct {
		name         string
		password     string
		passwordHash string
		wantErr      bool
	}{
		{
			name:         "Correct password",
			password:     password1,
			passwordHash: passwordHash1,
			wantErr:      false,
		},
		{
			name:         "Wrong password",
			password:     "wrongpass",
			passwordHash: passwordHash1,
			wantErr:      true,
		},
		{
			name:         "Empty password",
			password:     "",
			passwordHash: passwordHash1,
			wantErr:      true,
		},
		{
			name:         "Wrong hash",
			password:     password1,
			passwordHash: "wronghash",
			wantErr:      true,
		},
		{
			name:         "Another hash",
			password:     password1,
			passwordHash: passwordHash2,
			wantErr:      true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.passwordHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("err: %v, wantErr: %v", err, tt.wantErr)
			}
		})
	}
}

func TestMakeJWT(t *testing.T) {

	userID1 := uuid.New()
	tokenSecret1 := "asd123jdakldvaoijc"
	expiresIn := 1 * time.Second

	_, err := MakeJWT(userID1, tokenSecret1, expiresIn)
	if err != nil {
		t.Errorf("Shouldn't have error: %v", err)
	}
}

func TestValidateJWT(t *testing.T) {

	userID1 := uuid.New()
	tokenSecret := "asd123jdakldvaoijc"

	t.Run("Correct JWT", func(t *testing.T) {

		expiresIn := 1 * time.Second
		tokenString, _ := MakeJWT(userID1, tokenSecret, expiresIn)

		userID, err := ValidateJWT(tokenString, tokenSecret)
		if err != nil {
			t.Errorf("Shouldn't have error: %v", err)
		}
		if userID.String() != userID1.String() {
			t.Errorf("Expected user id: %v, but got: %v", userID1, userID)
		}
	})

	t.Run("Expired JWT", func(t *testing.T) {

		expiresIn := 1 * time.Second
		tokenString, _ := MakeJWT(userID1, tokenSecret, expiresIn)

		time.Sleep(1 * time.Second)
		_, err := ValidateJWT(tokenString, tokenSecret)
		if err == nil {
			t.Error("Should have error")
		}
	})

	t.Run("Wrong secret", func(t *testing.T) {
		expiresIn := 1 * time.Second
		tokenString, _ := MakeJWT(userID1, tokenSecret, expiresIn)

		wrongSecret := "wrongsecret123123"

		_, err := ValidateJWT(tokenString, wrongSecret)
		if err == nil {
			t.Error("Should have error")
		}
	})
}
