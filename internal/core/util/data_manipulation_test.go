package util

import (
	"testing"
)

func TestContains(t *testing.T) {
	type args struct {
		slice []int8
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "contains",
			args: args{
				slice: []int8{1, 2, 3, 4, 5},
			},
			wantErr: false,
		},
		{
			name: "not contains",
			args: args{
				slice: []int8{6, 7, 8, 9, 10},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isContains := Contains(tt.args.slice, 1, nil)
			if tt.wantErr && isContains {
				t.Errorf("Contains() = %v, want %v", isContains, tt.wantErr)
			}
			if !tt.wantErr && !isContains {
				t.Errorf("Contains() = %v, want %v", isContains, !tt.wantErr)
			}
		})
	}
}
