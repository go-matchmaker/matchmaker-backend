package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHashAndCompare(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "hash and compare successfully",
			args: args{
				password: "test",
			},
			wantErr: false,
		},
		{
			name: "hash too long",
			args: args{
				password: "01234567890123456789012345678901234567890123456789012345678901234567890123456789",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashed, err := HashPassword(tt.args.password)
			if tt.wantErr {
				assert.EqualError(t, err, bcrypt.ErrPasswordTooLong.Error())
				assert.Empty(t, hashed)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, hashed)

			err = ComparePassword(tt.args.password, hashed)
			assert.NoError(t, err)
		})
	}
}
