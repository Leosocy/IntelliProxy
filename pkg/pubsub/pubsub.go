// Copyright (c) 2019 leosocy, leosocy@gmail.com
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package pubsub

import (
	"sync"
)

// Notifier is an interface that send notification when data changes.
type Notifier interface {
	Attach(w Watcher)
	Detach(w Watcher)
	Notify(obj interface{})
}

// Watcher is an interface that receive notification that
// the notifier notify when the data changes.
type Watcher interface {
	// Receipt is called by Notifier.Notify.
	Receipt(obj interface{})
}

// BaseNotifier is a base implementation of the Notifier interface
type BaseNotifier struct {
	watchers []Watcher
	mu       sync.RWMutex
}

func (n *BaseNotifier) Attach(w Watcher) {
	if w == nil {
		panic("nil watcher")
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	n.watchers = append(n.watchers, w)
}

func (n *BaseNotifier) Detach(w Watcher) {
	if w == nil {
		panic("nil watcher")
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	for idx, v := range n.watchers {
		if v == w {
			copy(n.watchers[idx:], n.watchers[idx+1:])
			n.watchers = n.watchers[:len(n.watchers)-1]
			return
		}
	}
}

func (n *BaseNotifier) Notify(obj interface{}) {
	var wg sync.WaitGroup
	n.mu.RLock()
	wg.Add(len(n.watchers))
	for _, w := range n.watchers {
		go func(w Watcher) {
			defer wg.Done()
			w.Receipt(obj)
		}(w)
	}
	n.mu.RUnlock()
	wg.Wait()
}
