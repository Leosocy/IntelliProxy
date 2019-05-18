// Copyright (c) 2019 leosocy, leosocy@gmail.com
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package storage

import (
	"net"
	"testing"

	"github.com/Leosocy/IntelliProxy/pkg/proxy"
	"github.com/stretchr/testify/assert"
)

var testStorages = []Storage{
	NewInMemoryStorage(),
}

func TestInsert(t *testing.T) {
	assert := assert.New(t)
	for _, s := range testStorages {
		// insert invalid proxy
		err := s.Insert(nil)
		assert.NotNil(err)
		err = s.Insert(&proxy.Proxy{IP: net.ParseIP("1.2.3.4"), Port: 80, Score: 0})
		assert.NotNil(err)
		// insert one proxy
		err = s.Insert(&proxy.Proxy{IP: net.ParseIP("5.6.7.8"), Port: 80, Score: 100})
		assert.Nil(err)
		assert.Equal(uint(1), s.Len())
		// insert another proxy
		err = s.Insert(&proxy.Proxy{IP: net.ParseIP("9.10.11.12"), Port: 80, Score: 100})
		assert.Equal(uint(2), s.Len())
		// insert proxy with same IP, but diff score
		err = s.Insert(&proxy.Proxy{IP: net.ParseIP("9.10.11.12"), Port: 80, Score: 80})
		assert.Equal(uint(2), s.Len())
	}
}
