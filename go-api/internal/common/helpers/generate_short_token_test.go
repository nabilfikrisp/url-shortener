package helpers

import "testing"

func TestGenerateShortToken(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateShortToken(tt.args.s); got != tt.want {
				t.Errorf("GenerateShortToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
