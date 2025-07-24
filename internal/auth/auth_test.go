package auth

import (
	"testing"
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
