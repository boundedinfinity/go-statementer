package util

import (
	"fmt"

	"github.com/boundedinfinity/go-commoner/idiomatic/slicer"
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
)

func NewGenericErrorWrapper[T any](value T) *GenericErrorWrapper[T] {
	return &GenericErrorWrapper[T]{
		Value: value,
	}
}

type GenericErrorWrapper[T any] struct {
	UnwrapChain []error
	Value       T
	Message     string
}

func (this *GenericErrorWrapper[T]) WithErrs(errs ...error) *GenericErrorWrapper[T] {
	this.UnwrapChain = append(this.UnwrapChain, errs...)
	return this
}

func (this *GenericErrorWrapper[T]) WithMessage(message string) *GenericErrorWrapper[T] {
	this.Message = message
	return this
}

func (this *GenericErrorWrapper[T]) Error() string {
	var message []string

	message = append(message, slicer.JoinFn(
		func(_ int, err error) string { return err.Error() },
		" : ",
		slicer.Reverse(this.UnwrapChain...)...,
	))

	if this.Message != "" {
		message = append(message, this.Message)
	}

	message = append(message, fmt.Sprintf("%+v", this.Value))

	return stringer.Join(" : ", message...)
}

func (this *GenericErrorWrapper[T]) Unwrap() []error {
	return this.UnwrapChain
}
