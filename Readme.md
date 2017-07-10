# Deferred

A thread-safe [Deferred](https://developer.mozilla.org/en-US/docs/Mozilla/JavaScript_code_modules/Promise.jsm/Deferred#backwards_forwards_compatible) implementation for those times when you need to call `wg.Done()` multiple times in a `sync.WaitGroup` and not have it panic.

## Documentation

The full documentation is available on [Godoc](https://godoc.org/github.com/matthewmueller/go-deferred).

## Example

```go
ctx, cancel := context.WithCancel(context.Background())
d := deferred.New(ctx)
wg.Add(1)

go func() {
  v, err := d.Wait()
  if err != nil {
    t.Fatal(err)
  }
  assert.Equal(t, "a", v)
  wg.Done()
}()

d.Resolve("a")
d.Reject(errors.New("oh no")) // no-op
d.Resolve("hi") // no-op
v, err := d.Wait()
assert.Equal(t, "a", v) // always "a"
assert.Nil(t, err) // always nil
wg.Wait()
```

## License

MIT
