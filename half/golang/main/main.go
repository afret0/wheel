package main

import (
	"fmt"
)

func main() {
	// ctx := context.Background()
	// ctx = context.WithValue(ctx, "key", "value")
	// ctx = context.WithValue(ctx, "key2", false)

	// key := ctx.Value("key").(string)
	// key, ok := ctx.Value("key").(string)
	// fmt.Println(key, ok)

	// // key1 := ctx.Value("key1").(string)
	// key1, ok := ctx.Value("key1").(string)
	// fmt.Println(key1, ok)

	// key2 := ctx.Value("key2").(bool)
	// key2, ok = ctx.Value("key2").(bool)
	// fmt.Println(key2, ok)

	m1 := map[string]int{"a": 1}
	key3, ok := m1["a"]
	fmt.Println(key3, ok)
	key4, ok := m1["key4"]
	fmt.Println(key4, ok)
	key5 :=m1["key5"]
	fmt.Println(key5)
}
