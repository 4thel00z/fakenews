package fakenews

import (
	"math/rand"
	"reflect"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Runner[T any] func(current *T) (T, error)
type Rule[T any] interface {
	Run(current *T) (T, error)
	FieldName() string
}

type OneFieldRule[T any] struct {
	fieldName string
	runner    Runner[T]
}

func NewRule[T any](fn string, runner Runner[T]) OneFieldRule[T] {
	return OneFieldRule[T]{
		fieldName: fn,
		runner:    runner,
	}
}

func (o OneFieldRule[T]) Run(current *T) (T, error) {
	return o.runner(current)
}

func (o OneFieldRule[T]) FieldName() string {
	return o.fieldName
}

// RandomInt is a generator that will set a random integer between the given start and end values
// to the field.
func RandomInt[T any](fieldName string, start, end int) Rule[T] {
	return NewRule(fieldName, func(current *T) (_ T, err error) {
		defer func() {
			if r := recover(); r != nil {
				err, _ = r.(error)
			}
		}()
		reflect.ValueOf(current).Elem().FieldByName(fieldName).Set(reflect.ValueOf(rand.Intn(end-start) + start))
		return *current, err
	})
}

// RandomFloat is a generator that will set a random float between the given start and end values
// to the field.
func RandomFloat[T any](fieldName string, start, end float64) Rule[T] {
	return NewRule(fieldName, func(current *T) (_ T, err error) {
		defer func() {
			if r := recover(); r != nil {
				err, _ = r.(error)
			}
		}()
		reflect.ValueOf(current).Elem().FieldByName(fieldName).Set(reflect.ValueOf((rand.Float64() * (end - start)) + start))
		return *current, err
	})
}

// RandomPick is a generator that will pick a random value from the given list of values
// and set it to the field.
func RandomPick[T any, F any](fieldName string, vals ...F) Rule[T] {
	return NewRule(fieldName, func(current *T) (_ T, err error) {
		defer func() {
			if r := recover(); r != nil {
				err, _ = r.(error)
			}
		}()
		reflect.ValueOf(current).Elem().FieldByName(fieldName).Set(reflect.ValueOf(vals[rand.Intn(len(vals))]))
		return *current, err
	})
}

// ArrayGenerator is a generator that will set a random array of values from the given list of generators
// to the field.
func ArrayGenerator[T any](fieldName string, n int, gs ...Generator[T]) Rule[T] {
	return NewRule(fieldName, func(current *T) (_ T, err error) {
		defer func() {
			if r := recover(); r != nil {
				err, _ = r.(error)
			}
		}()
		reflect.ValueOf(current).Elem().FieldByName(fieldName).Set(reflect.ValueOf(ConsumeUntil(Infinite(ChainRandom(n, gs...)), n)))
		return *current, err
	})
}

// RoundRobinGenerator is a generator that will iterate over a list of generators in a round robin fashion
// until the limit is reached or one of the generators is exhausted.
func RoundRobinGenerator[T any](fieldName string, n int, gs ...Generator[T]) Rule[T] {
	return NewRule(fieldName, func(current *T) (_ T, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = r.(error)
			}
		}()
		reflect.ValueOf(current).Elem().FieldByName(fieldName).Set(reflect.ValueOf(ConsumeUntil(ChainRoundRobin(n, gs...), n)))
		return *current, err
	})
}

func RulesToGenerator[T any](t T, rs ...Rule[T]) Generator[T] {
	s := StepGenerator[T]{}
	for _, r := range rs {
		s.Add(func() T {
			run, _ := r.Run(&t)
			return run
		})
	}

	return s
}
