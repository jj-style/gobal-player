package resty_test

import (
	"errors"
	"testing"

	"github.com/jj-style/gobal-player/pkg/resty"
	"github.com/jj-style/gobal-player/pkg/resty/mocks"
	"github.com/stretchr/testify/assert"
)

type testType struct {
	String string `json:"string"`
	Int    int    `json:"int"`
}

func TestGet(t *testing.T) {
	type deps struct {
		client *mocks.MockClient
	}
	tests := []struct {
		name     string
		setup    func(deps)
		want     testType
		checkErr assert.ErrorAssertionFunc
	}{
		{
			name: "happy",
			setup: func(d deps) {
				d.client.EXPECT().Get("url").Return([]byte(`{"string": "abc", "int": 123}`), nil)
			},
			want: testType{
				String: "abc",
				Int:    123,
			},
			checkErr: assert.NoError,
		},
		{
			name: "unhappy request",
			setup: func(d deps) {
				d.client.EXPECT().Get("url").Return([]byte{}, errors.New("boom"))
			},
			checkErr: assert.Error,
		},
		{
			name: "unhappy json",
			setup: func(d deps) {
				d.client.EXPECT().Get("url").Return([]byte("i am not json"), nil)
			},
			checkErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := deps{
				client: mocks.NewMockClient(t),
			}
			if tt.setup != nil {
				tt.setup(deps)
			}

			got, err := resty.Get[testType](deps.client, "url")
			tt.checkErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPost(t *testing.T) {
	type deps struct {
		client *mocks.MockClient
	}
	type args struct {
		req testType
	}
	tests := []struct {
		name     string
		setup    func(deps)
		args     args
		want     testType
		checkErr assert.ErrorAssertionFunc
	}{
		{
			name: "happy",
			setup: func(d deps) {
				d.client.EXPECT().Post("url", testType{String: "abc", Int: 123}).Return([]byte(`{"string": "xyz", "int": 789}`), nil)
			},
			args: args{
				req: testType{
					String: "abc",
					Int:    123,
				},
			},
			want: testType{
				String: "xyz",
				Int:    789,
			},
			checkErr: assert.NoError,
		},
		{
			name: "unhappy request",
			setup: func(d deps) {
				d.client.EXPECT().Post("url", testType{String: "abc", Int: 123}).Return([]byte{}, errors.New("boom"))
			},
			args: args{
				req: testType{
					String: "abc",
					Int:    123,
				},
			},
			checkErr: assert.Error,
		},
		{
			name: "unhappy json",
			setup: func(d deps) {
				d.client.EXPECT().Post("url", testType{String: "abc", Int: 123}).Return([]byte("i am not json"), nil)
			},
			args: args{
				req: testType{
					String: "abc",
					Int:    123,
				},
			},
			checkErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := deps{
				client: mocks.NewMockClient(t),
			}
			if tt.setup != nil {
				tt.setup(deps)
			}

			got, err := resty.Post[testType, testType](deps.client, "url", tt.args.req)
			tt.checkErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
