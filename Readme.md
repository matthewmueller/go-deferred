# Deferred

A thread-safe [Deferred](https://developer.mozilla.org/en-US/docs/Mozilla/JavaScript_code_modules/Promise.jsm/Deferred#backwards_forwards_compatible) implementation for those times when you need to call `wg.Done()` multiple times in a `sync.WaitGroup` and not have it panic.

`deferred.Wait()` cancellable with a context.

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
wg.Wait()
```

## License

MIT
