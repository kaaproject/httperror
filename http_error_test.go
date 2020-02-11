// Copyright 2019 KaaIoT Technologies, LLC

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httperror

import (
	"errors"
	"net/http"
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
		{name: "200", fields: fields{statusCode: http.StatusOK, description: ""}, want: ""},
		{name: "404", fields: fields{statusCode: http.StatusNotFound, description: "no such page"}, want: "no such page"},
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
		{
			name: "200",
			args: args{code: http.StatusOK, format: "", a: nil},
			want: &HTTPError{statusCode: http.StatusOK, description: ""},
		},
		{
			name: "404",
			args: args{code: http.StatusNotFound, format: "no such page", a: nil},
			want: &HTTPError{statusCode: http.StatusNotFound, description: "no such page"},
		},
		{
			name: "429",
			args: args{code: http.StatusTooManyRequests, format: "over %v requests", a: a},
			want: &HTTPError{statusCode: http.StatusTooManyRequests, description: "over 9000 requests"},
		},
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
		{name: "nil", args: args{err: nil}, want: http.StatusOK},
		{name: "regular error", args: args{err: errors.New("any regular error")}, want: http.StatusInternalServerError},
		{name: "http error", args: args{err: New(http.StatusNotFound, "no such page")}, want: http.StatusNotFound},
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
		{name: "http error", args: args{err: New(http.StatusNotFound, "no such page")}, want: "no such page"},
		{name: "zero http error", args: args{err: New(0, "")}, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusText(tt.args.err); got != tt.want {
				t.Errorf("StatusText() = %v, want %v", got, tt.want)
			}
		})
	}
}
