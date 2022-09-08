package examples

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/si3nloong/sqlike/sql/codec"
)

func testRace(ctx context.Context, t *testing.T) {
	registry := codec.DefaultRegistry
	wg := new(sync.WaitGroup)
	getStruct := func(v any) {
		defer wg.Done()
		to := reflect.TypeOf(v)
		if _, err := registry.LookupDecoder(to); err != nil {
			panic(err)
		}
	}
	go getStruct(struct {
		Name string
		Age  int
	}{})
	go getStruct(struct {
		Name string
		Age  int
	}{})
	go getStruct(struct {
		Name string
		Age  int
	}{})
	wg.Add(3)
	wg.Wait()
}
