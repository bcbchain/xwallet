package client_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/rpc/client/mock"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func TestWaitForHeight(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	m := &mock.StatusMock{
		Call: mock.Call{
			Error: errors.New("bye"),
		},
	}
	r := mock.NewStatusRecorder(m)

	err := client.WaitForHeight(r, 8, nil)
	require.NotNil(err)
	require.Equal("bye", err.Error())

	require.Equal(1, len(r.Calls))

	m.Call = mock.Call{
		Response: &ctypes.ResultStatus{SyncInfo: ctypes.SyncInfo{LatestBlockHeight: 10}},
	}

	err = client.WaitForHeight(r, 40, nil)
	require.NotNil(err)
	require.True(strings.Contains(err.Error(), "aborting"))

	require.Equal(2, len(r.Calls))

	err = client.WaitForHeight(r, 5, nil)
	require.Nil(err)

	require.Equal(3, len(r.Calls))

	myWaiter := func(delta int64) error {

		m.Call.Response = &ctypes.ResultStatus{SyncInfo: ctypes.SyncInfo{LatestBlockHeight: 15}}
		return client.DefaultWaitStrategy(delta)
	}

	err = client.WaitForHeight(r, 12, myWaiter)
	require.Nil(err)

	require.Equal(5, len(r.Calls))

	pre := r.Calls[3]
	require.Nil(pre.Error)
	prer, ok := pre.Response.(*ctypes.ResultStatus)
	require.True(ok)
	assert.Equal(int64(10), prer.SyncInfo.LatestBlockHeight)

	post := r.Calls[4]
	require.Nil(post.Error)
	postr, ok := post.Response.(*ctypes.ResultStatus)
	require.True(ok)
	assert.Equal(int64(15), postr.SyncInfo.LatestBlockHeight)
}
