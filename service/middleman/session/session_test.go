// Copyright (c) 2019 leosocy, leosocy@gmail.com
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package session

import (
	"testing"

	"github.com/Leosocy/IntelliProxy/pkg/proxy"
	"github.com/stretchr/testify/assert"
)

func TestSessionPoolManipulate(t *testing.T) {
	assert := assert.New(t)
	pool := &sessionPool{}
	// new two sessions
	pxyOne, _ := proxy.NewProxy("1.2.3.4", "8888")
	pxyAnother, _ := proxy.NewProxy("5.6.7.8", "9999")
	one := pool.new(pxyOne, newDefaultSessionTransport())
	another := pool.new(pxyAnother, newDefaultSessionTransport())
	assert.NotNil(one)
	assert.NotNil(another)
	// put all
	pool.put(one)
	pool.put(one)
	pool.put(another)
	assert.Equal(2, len(pool.sessions))
	// remove
	pool.remove(one)
	assert.NotContains(pool.sessions, one)
	pool.remove(pool.new(pxyOne, newDefaultSessionTransport()))
	assert.NotEmpty(pool.sessions)
	pool.remove(another)
	assert.NotContains(pool.sessions, another)
	assert.Empty(pool.sessions)
}
