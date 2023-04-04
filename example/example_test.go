package example

import (
	"errors"
	"testing"
)

func Test_add(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1 + 2 = 3",
			args: args{
				a: 1,
				b: 2,
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := add(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_divide(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1 / 2 = 0",
			args: args{
				a: 1,
				b: 2,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := divide(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("divide() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_multiply(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1 * 2 = 2",
			args: args{
				a: 1,
				b: 2,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := multiply(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("multiply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_multiplyString(t *testing.T) {
	type args struct {
		a int
		b string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "2 * two = twotwo",
			args: args{
				a: 2,
				b: "two",
			},
			want: "twotwo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := multiplyString(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("multiplyString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_subtract(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1 - 2 = -1",
			args: args{
				a: 1,
				b: 2,
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := subtract(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("subtract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertRomanToDecimal(t *testing.T) {
	testCases := []struct {
		roman    string
		expected int
		err      error
	}{
		{"III", 3, nil},
		{"IV", 4, nil},
		{"IX", 9, nil},
		{"LVIII", 58, nil},
		{"MCMXCIV", 1994, nil},
		{"IIII", 0, errors.New("invalid Roman numeral")},
		{"VV", 0, errors.New("invalid Roman numeral")},
	}

	for _, tc := range testCases {
		t.Run(tc.roman, func(t *testing.T) {
			result, err := ConvertRomanToDecimal(tc.roman)
			if err != nil && tc.err == nil {
				t.Errorf("Unexpected error for Roman numeral %s: %v", tc.roman, err)
			} else if err == nil && tc.err != nil {
				t.Errorf("Expected error but got none for Roman numeral %s", tc.roman)
			} else if err != nil && tc.err != nil && err.Error() != tc.err.Error() {
				t.Errorf("Expected error %v but got %v for Roman numeral %s", tc.err, err, tc.roman)
			}

			if result != tc.expected {
				t.Errorf("Expected %d but got %d for Roman numeral %s", tc.expected, result, tc.roman)
			}
		})
	}
}
