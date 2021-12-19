package framework

import (
	"golang.org/x/sync/errgroup"
	"strconv"
	"strings"
	"testing"
	"time"
)

type FinalizeFunc func() error

type TestContext struct {
	Id           string
	cleanupFuncs []FinalizeFunc // usually resource delete functions
}

func (f *Framework) NewTestContext(t *testing.T) TestContext {
	prefix := strings.TrimPrefix(
		strings.Replace(strings.ToLower(t.Name()), "/", "-", -1),
		"test",
	)
	id := prefix + "-" + strconv.FormatInt(time.Now().Unix(), 36)
	return TestContext{Id: id}
}

// GetObjId returns an ascending ID based on the length of cleanUpFns. It is
// based on the premise that every new object also appends a new finalizerFn on
// cleanUpFns. This can e.g. be used to create multiple namespaces in the same
// test context.
func (ctx *TestContext) GetObjId() string {
	return ctx.Id + "-" + strconv.Itoa(len(ctx.cleanupFuncs))
}

func (ctx *TestContext) Cleanup(t *testing.T) {
	var eg errgroup.Group
	for i := len(ctx.cleanupFuncs) - 1; i >= 0; i-- {
		// call the given func in a new goroutine
		eg.Go(ctx.cleanupFuncs[i])
	}
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}

func (ctx *TestContext) AddFinalizerFunc(fn FinalizeFunc) {
	ctx.cleanupFuncs = append(ctx.cleanupFuncs, fn)
}
