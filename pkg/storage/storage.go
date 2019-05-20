// Copyright (c) 2019 leosocy, leosocy@gmail.com
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package storage

import (
	"errors"
	"net"

	"github.com/Leosocy/IntelliProxy/pkg/proxy"
)

var (
	ErrProxyInvalid       = errors.New("proxy is nil or score <= 0")
	ErrProxyDuplicated    = errors.New("proxy is already in storage")
	ErrProxyDoesNotExists = errors.New("proxy doesn't exists")
)

type QueryCondition struct {
}

type Storage interface {
	Insert(p *proxy.Proxy) error
	Update(newP *proxy.Proxy) error
	Search(ip net.IP) *proxy.Proxy
	Delete(ip net.IP) error
	Best() *proxy.Proxy
	Len() uint
	// Query(cond QueryCondition) ([]*proxy.Proxy, error)
}
