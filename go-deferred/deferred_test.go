package deferred_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/matthewmueller/deferred"
	"github.com/stretchr/testify/assert"
)

func TestResolve(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	d := deferred.New(ctx)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		v, err := d.Wait()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "a", v)
		wg.Done()
	}()

	go func() {
		v, err := d.Wait()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "a", v)
		wg.Done()
	}()

	d.Resolve("a")
	d.Resolve("b")
	d.Resolve("c")
	d.Resolve("e")
	wg.Wait()
	cancel()
	v, err := d.Wait()
	assert.Equal(t, "a", v)
	assert.Nil(t, err)
	// assert.EqualError(t, err, "context canceled")
}

func TestReject(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	d := deferred.New(ctx)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		v, err := d.Wait()
		assert.Nil(t, v)
		assert.EqualError(t, err, "oh no")
		wg.Done()
	}()

	go func() {
		v, err := d.Wait()
		assert.EqualError(t, err, "oh no")
		assert.Nil(t, v)
		wg.Done()
	}()

	d.Reject(errors.New("oh no"))
	cancel()
	d.Resolve("a")
	d.Resolve("b")
	d.Resolve("c")
	d.Resolve("e")
	wg.Wait()

	v, err := d.Wait()
	assert.Nil(t, v)
	assert.EqualError(t, err, "oh no")
}
