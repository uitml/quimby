package validate

import (
	"testing"
)

func TestDefaultIfEmpty(t *testing.T) {
	type args struct {
		field         string
		default_value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// Testcase 1: field empty (default_value returned)
		{
			name: "Field empty, default returned",
			args: args{
				field:         "",
				default_value: "bar",
			},
			want: "bar",
		},
		// Testcase 2: field not empty (field returned)
		{
			name: "Field filled, field returned",
			args: args{
				field:         "foo",
				default_value: "bar",
			},
			want: "foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultIfEmpty(tt.args.field, tt.args.default_value); got != tt.want {
				t.Errorf("DefaultIfEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsername(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// Testcase 1: valid username foo123. Return true
		{
			name: "valid username",
			args: args{"foo123"},
			want: true,
		},
		// Testcase 2: invalid (too many letters at start) afoo123. Return false
		{
			name: "too many letters",
			args: args{"afoo123"},
			want: false,
		},
		// Testcase 3: invalid (too many numbers at end) foo1234. Return false
		{
			name: "too many numbers",
			args: args{"foo1234"},
			want: false,
		},
		// Testcase 4: invalid (number at start + valid username) 1foo123. Return false
		{
			name: "number before valid",
			args: args{"1foo123"},
			want: false,
		},
		// Testcase 5: invalid (valid username + letter at end) foo123a. Return false
		{
			name: "letter after valid",
			args: args{"foo123a"},
			want: false,
		},
		// Testcase 6: invalid (number at start) 1oo123. Return false
		{
			name: "number at start",
			args: args{"1oo123"},
			want: false,
		},
		// Testcase 7: invalid (letter at end) foo12c. Return false
		{
			name: "number at start",
			args: args{"foo12c"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Username(tt.args.username); got != tt.want {
				t.Errorf("Username() = %v, want %v", got, tt.want)
			}
		})
	}
}
