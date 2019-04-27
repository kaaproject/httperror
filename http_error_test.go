package httperror

import (
	"errors"
	"reflect"
	"testing"
)

func TestHTTPError_Error(t *testing.T) {
	type fields struct {
		statusCode  int
		description string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "200", fields: fields{statusCode: 200, description: ""}, want: ""},
		{name: "404", fields: fields{statusCode: 404, description: "no such page"}, want: "no such page"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HTTPError{
				statusCode:  tt.fields.statusCode,
				description: tt.fields.description,
			}
			if got := p.Error(); got != tt.want {
				t.Errorf("HTTPError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		code   int
		format string
		a      []interface{}
	}
	a := make([]interface{}, 1)
	a[0] = 9000

	tests := []struct {
		name string
		args args
		want *HTTPError
	}{
		{name: "200", args: args{code: 200, format: "", a: nil}, want: &HTTPError{statusCode: 200, description: ""}},
		{name: "404", args: args{code: 404, format: "no such page", a: nil}, want: &HTTPError{statusCode: 404, description: "no such page"}},
		{name: "429", args: args{code: 429, format: "%v requests in last minute", a: a}, want: &HTTPError{statusCode: 429, description: "9000 requests in last minute"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.code, tt.args.format, tt.args.a...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusCode(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "nil", args: args{err: nil}, want: 200},
		{name: "regular error", args: args{err: errors.New("any regular error")}, want: 500},
		{name: "http error", args: args{err: New(404, "no such page")}, want: 404},
		{name: "zero http error", args: args{err: New(0, "")}, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusCode(tt.args.err); got != tt.want {
				t.Errorf("StatusCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReasonPhrase(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "nil", args: args{err: nil}, want: "OK"},
		{name: "regular error", args: args{err: errors.New("any regular error")}, want: "Internal Server Error"},
		{name: "http error", args: args{err: New(404, "no such page")}, want: "Not Found"},
		{name: "zero http error", args: args{err: New(0, "")}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReasonPhrase(tt.args.err); got != tt.want {
				t.Errorf("ReasonPhrase() = %v, want %v", got, tt.want)
			}
		})
	}
}
