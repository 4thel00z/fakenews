package fakenews

import (
	"math/rand"
	"reflect"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Runner func(current interface{}) (interface{}, error)
type Rule interface {
	Run(current interface{}) (interface{}, error)
	FieldName() string
}

type OneFieldRule struct {
	fieldName string
	runner    Runner
}

func NewRule(fn string, runner Runner) OneFieldRule {
	return OneFieldRule{
		fieldName: fn,
		runner:    runner,
	}
}

func (o OneFieldRule) Run(current interface{}) (interface{}, error) {
	return o.runner(current)
}

func (o OneFieldRule) FieldName() string {
	return o.fieldName
}

func RandomInt(fieldName string, start, end int) Rule {
	return NewRule(fieldName, func(current interface{}) (interface{}, error) {
		reflect.ValueOf(&current).Elem().FieldByName(fieldName).Set(reflect.ValueOf(rand.Intn(end-start+1) + start))
		return current, nil
	})
}

func RandomFloat(fieldName string, start, end float64) Rule {
	return NewRule(fieldName, func(current interface{}) (interface{}, error) {
		reflect.ValueOf(&current).Elem().FieldByName(fieldName).Set(reflect.ValueOf((rand.Float64() * (end - start)) + start))
		return current, nil
	})
}

func RandomPick(fieldName string, vals ...interface{}) Rule {
	return NewRule(fieldName, func(current interface{}) (interface{}, error) {
		reflect.ValueOf(&current).Elem().FieldByName(fieldName).Set(reflect.ValueOf(vals[rand.Intn(len(vals))]))
		return current, nil
	})
}

func RandomGenerator(fieldName string, n int, gs ...Generator) Rule {
	return NewRule(fieldName, func(current interface{}) (interface{}, error) {
		reflect.ValueOf(&current).Elem().FieldByName(fieldName).Set(reflect.ValueOf(ConsumeUntil(n, ChainRandom(n, gs...))))
		return current, nil
	})
}

func RoundRobinGenerator(fieldName string, n int, gs ...Generator) Rule {
	return NewRule(fieldName, func(current interface{}) (interface{}, error) {
		reflect.ValueOf(&current).Elem().FieldByName(fieldName).Set(reflect.ValueOf(ConsumeUntil(n, ChainRoundRobin(n, gs...))))
		return current, nil
	})
}
