package fakenews

import (
	"bufio"
	"io"
	"math/rand"
	"sync"
)

type Generator[T any] interface {
	Next() (T, bool)
	Iter() Generator[T]
}

func GeneratorFromIO[T any](file io.ReadSeeker, parser LineParser[T]) Generator[T] {
	return LineGenerator[T]{
		file:   file,
		parser: parser,
		mx:     &sync.Mutex{},
		reader: bufio.NewReader(file),
	}
}

type LineGenerator[T any] struct {
	file   io.ReadSeeker
	parser LineParser[T]
	reader *bufio.Reader
	mx     *sync.Mutex
}

func (l LineGenerator[T]) Next() (T, bool) {
	l.mx.Lock()
	defer l.mx.Unlock()

	var t T
	if l.reader == nil {
		l.reader = bufio.NewReader(l.file)
	}
	line, _, err := l.reader.ReadLine()
	if err != nil {
		return t, false
	}

	if l.parser == nil {
		return t, false
	}

	return l.parser(line), true
}

func (l LineGenerator[T]) Iter() Generator[T] {
	l.file.Seek(0, 0)
	return LineGenerator[T]{
		file:   l.file,
		parser: l.parser,
		mx:     l.mx,
		reader: bufio.NewReader(l.file),
	}
}

type LineParser[T any] func(line []byte) T

type Step[T any] func() T

type StepGenerator[T any] struct {
	current int
	steps   []Step[T]
}

func (s StepGenerator[T]) Next() (T, bool) {
	if s.current >= len(s.steps)-1 {
		var t T
		return t, false
	}
	return s.steps[s.current](), true
}

func (s StepGenerator[T]) Iter() Generator[T] {
	return StepGenerator[T]{
		current: 0,
		steps:   s.steps,
	}
}

func (sg StepGenerator[T]) Add(s Step[T]) StepGenerator[T] {
	return StepGenerator[T]{
		current: sg.current,
		steps:   append(sg.steps, s),
	}
}

func ConsumeGenerator[T any](g Generator[T]) []T {
	vals := []T{}

	for {
		val, done := g.Next()
		if done {
			break
		}
		vals = append(vals, val)
	}
	return vals
}

func ConsumeUntil[T any](g Generator[T], n int) []T {
	vals := []T{}

	for n > 0 {
		val, done := g.Next()
		if done {
			break
		}
		vals = append(vals, val)
		n--
	}
	return vals
}

type InfiniteGenerator[T any] struct {
	generator Generator[T]
}

func (i InfiniteGenerator[T]) Next() (T, bool) {
	next, b := i.generator.Next()
	if !b {
		i.generator = i.generator.Iter()
		return i.generator.Next()
	}
	return next, true
}

func (i InfiniteGenerator[T]) Iter() Generator[T] {
	return InfiniteGenerator[T]{
		generator: i.generator.Iter(),
	}
}

type RoundRobinChainGenerator[T any] struct {
	current       int
	generators    []Generator[T]
	limit         int
	originalLimit int
}

func (c RoundRobinChainGenerator[T]) Next() (T, bool) {
	generator := c.generators[c.current]

	nextVal, ok := generator.Next()
	if !ok {
		c.generators[c.current] = generator.Iter()
		return c.Next()
	}

	c.current++
	c.current %= len(c.generators)
	c.limit--

	return nextVal, c.limit >= 0
}

func (c RoundRobinChainGenerator[T]) Iter() Generator[T] {
	return RoundRobinChainGenerator[T]{
		current:       c.current,
		generators:    c.generators,
		limit:         c.originalLimit,
		originalLimit: c.originalLimit,
	}
}

type RandomChainGenerator[T any] struct {
	generators    []Generator[T]
	limit         int
	originalLimit int
}

func (c RandomChainGenerator[T]) Next() (T, bool) {
	current := rand.Int() % len(c.generators)
	generator := c.generators[current]

	nextVal, ok := generator.Next()
	if !ok {
		c.generators[current] = generator.Iter()
		return c.Next()
	}

	c.limit--

	return nextVal, c.limit >= 0
}

func (c RandomChainGenerator[T]) Iter() Generator[T] {
	return RandomChainGenerator[T]{
		generators:    c.generators,
		limit:         c.originalLimit,
		originalLimit: c.originalLimit,
	}
}

func ChainRoundRobin[T any](limit int, g ...Generator[T]) Generator[T] {
	return RoundRobinChainGenerator[T]{
		current:       0,
		generators:    g,
		limit:         limit,
		originalLimit: limit,
	}
}
func ChainRandom[T any](limit int, g ...Generator[T]) Generator[T] {
	return RandomChainGenerator[T]{
		generators:    g,
		limit:         limit,
		originalLimit: limit,
	}
}

func Infinite[T any](g Generator[T]) Generator[T] {
	return InfiniteGenerator[T]{
		generator: g,
	}
}
