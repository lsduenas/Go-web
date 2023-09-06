package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()

	newCtx := addValue(ctx)

	s := newCtx.Value("bootcamp")

	fmt.Println(s)
}

func addValue(ctx context.Context) context.Context {
	return context.WithValue(ctx, "bootcamp", "GO")
}
