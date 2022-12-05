package main

import (
	"reflect"
	"testing"
)

func TestStack(t *testing.T) {
	stack := stack{}

	stack.pushTop("a")
	stack.pushTop("a")
	stack.pushTop("b")
	stack.pushBot("c")

	reflect.DeepEqual(stack, []string{"c", "a", "a", "b"})
}
