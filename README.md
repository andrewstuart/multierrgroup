Multierrgroup provides a marriage of
[x/sync/errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup.Group) and
hashicorp's [multierror](https://pkg.go.dev/github.com/hashicorp/go-multierror).
Start up as many error-returning goroutines as you need and get a multierror
with all the failures that occurred once they're all finished.

```go var meg multerrgroup.Group

meg.Go(func() error { // Do something dangerous })

meg.Go(func() error { // Try something you know will probably never work })

err := meg.Wait() ```
