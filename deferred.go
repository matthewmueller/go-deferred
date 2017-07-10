package deferred

import (
	"context"
	"sync"
)

// Deferred struct
type Deferred struct {
	ctx      context.Context
	once     *sync.Once
	resolved chan bool
	rejected chan bool
	value    interface{}
	err      error
}

// New Deferred
func New(ctx context.Context) *Deferred {
	return &Deferred{
		once:     &sync.Once{},
		resolved: make(chan bool),
		rejected: make(chan bool),
		ctx:      ctx,
	}
}

// Resolve the deferred
func (d *Deferred) Resolve(v interface{}) {
	d.once.Do(func() {
		d.value = v
		close(d.resolved)
	})
}

// Reject the deferred
func (d *Deferred) Reject(err error) {
	d.once.Do(func() {
		d.err = err
		close(d.rejected)
	})
}

// Wait returns a value if the deferred was resolved or
// an error if the deferred was rejected or cancelled
func (d *Deferred) Wait() (v interface{}, err error) {
	select {
	case <-d.resolved:
		return d.value, nil
	case <-d.rejected:
		return nil, d.err
	case <-d.ctx.Done():
		// this will take precedent over closed
		// channels, so we need to check if
		// any of the channels have been
		// already closed first
		select {
		case <-d.resolved:
			return d.value, nil
		case <-d.rejected:
			return nil, d.err
		default:
			return nil, d.ctx.Err()
		}
	}
}
