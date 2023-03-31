package example

import "testing"

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
