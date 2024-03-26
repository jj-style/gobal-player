package resty

import (
	"context"
	"errors"
	"testing"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNilCache(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)
	nc := &nilCache[int]{}

	var got int
	var err error

	// get
	got, err = nc.Get(ctx, "a")
	assert.NoError(err)
	assert.Equal(0, got)

	// set
	err = nc.Set(ctx, "a", 1)
	assert.NoError(err)

	// get still nil
	got, err = nc.Get(ctx, "a")
	assert.NoError(err)
	assert.Equal(0, got)

	// clear cache
	err = nc.Clear(ctx)
	assert.NoError(err)
}

func TestCache_Get(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		setup     func(s *store.MockStoreInterface)
		want      int
		assertErr assert.ErrorAssertionFunc
	}{
		{
			name: "cache hit",
			key:  "key",
			setup: func(s *store.MockStoreInterface) {
				s.EXPECT().Get(gomock.Any(), "key").Return(1, nil)
			},
			want:      1,
			assertErr: assert.NoError,
		},
		{
			name: "cache miss",
			key:  "key",
			setup: func(s *store.MockStoreInterface) {
				s.EXPECT().Get(gomock.Any(), "key").Return(0, errors.New("item missing"))
			},
			want:      0,
			assertErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			gomockController := gomock.NewController(t)
			mstore := store.NewMockStoreInterface(gomockController)

			if tt.setup != nil {
				tt.setup(mstore)
			}

			c := &cacheImpl[int]{cache: cache.New[int](mstore)}

			var got int
			var err error

			got, err = c.Get(ctx, tt.key)
			tt.assertErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCache_Set(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     int
		setup     func(s *store.MockStoreInterface)
		assertErr assert.ErrorAssertionFunc
	}{
		{
			name:  "happy",
			key:   "key",
			value: 1,
			setup: func(s *store.MockStoreInterface) {
				s.EXPECT().Set(gomock.Any(), "key", 1).Return(nil)
			},
			assertErr: assert.NoError,
		},
		{
			name:  "unhappy",
			key:   "key",
			value: 1,
			setup: func(s *store.MockStoreInterface) {
				s.EXPECT().Set(gomock.Any(), "key", 1).Return(errors.New("boom"))
			},
			assertErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			gomockController := gomock.NewController(t)
			mstore := store.NewMockStoreInterface(gomockController)

			if tt.setup != nil {
				tt.setup(mstore)
			}

			c := &cacheImpl[int]{cache: cache.New[int](mstore)}

			err := c.Set(ctx, tt.key, tt.value)
			tt.assertErr(t, err)
		})
	}
}

func TestCache_Clear(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(s *store.MockStoreInterface)
		assertErr assert.ErrorAssertionFunc
	}{
		{
			name: "happy",
			setup: func(s *store.MockStoreInterface) {
				s.EXPECT().Clear(gomock.Any()).Return(nil)
			},
			assertErr: assert.NoError,
		},
		{
			name: "unhappy",
			setup: func(s *store.MockStoreInterface) {
				s.EXPECT().Clear(gomock.Any()).Return(errors.New("boom"))
			},
			assertErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			gomockController := gomock.NewController(t)
			mstore := store.NewMockStoreInterface(gomockController)

			if tt.setup != nil {
				tt.setup(mstore)
			}

			c := &cacheImpl[int]{cache: cache.New[int](mstore)}

			err := c.Clear(ctx)
			tt.assertErr(t, err)
		})
	}
}
