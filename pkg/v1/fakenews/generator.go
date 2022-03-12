package fakenews

import (
	"bufio"
	"io"
	"math/rand"
)

type Generator interface {
	Next() (interface{}, bool)
	Iter() Generator
}

func GeneratorFromIO(file io.ReadSeeker, parser LineParser) Generator {
	return LineGenerator{
		file:   file,
		parser: parser,
	}
}

// Crucial note: LineGenerator is NOT threadsafe !
type LineGenerator struct {
	file   io.ReadSeeker
	parser LineParser
}

func (l LineGenerator) Next() (interface{}, bool) {
	line, _, err := bufio.NewReader(l.file).ReadLine()
	if err != nil {
		return "", false
	}

	if l.parser != nil {
		return l.parser(line), true
	}

	return line, true
}

func (l LineGenerator) Iter() Generator {
	l.file.Seek(0, 0)
	return LineGenerator{
		file:   l.file,
		parser: l.parser,
	}
}

type LineParser func(line []byte) interface{}

type Step func() interface{}

type StepGenerator struct {
	current int
	steps   []Step
}

func (s StepGenerator) Next() (interface{}, bool) {
	if s.current >= len(s.steps)-1 {
		return nil, false
	}
	return s.steps[s.current], true
}

func (s StepGenerator) Iter() Generator {
	return StepGenerator{
		current: 0,
		steps:   s.steps,
	}
}

func (sg StepGenerator) Add(s Step) StepGenerator {
	return StepGenerator{
		current: sg.current,
		steps:   append(sg.steps, s),
	}
}

func ConsumeGenerator(g Generator) []interface{} {
	vals := []interface{}{}

	for {
		val, done := g.Next()
		if done {
			break
		}
		vals = append(vals, val)
	}
	return vals
}
func ConsumeUntil(n int, g Generator) []interface{} {
	vals := []interface{}{}

	for {
		if n <= 0 {
			break
		}
		val, done := g.Next()
		if done {
			break
		}
		vals = append(vals, val)
		n--
	}
	return vals
}

type RoundRobinChainGenerator struct {
	current       int
	generators    []Generator
	limit         int
	originalLimit int
}

func (c RoundRobinChainGenerator) Next() (interface{}, bool) {
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

func (c RoundRobinChainGenerator) Iter() Generator {
	return RoundRobinChainGenerator{
		current:       c.current,
		generators:    c.generators,
		limit:         c.originalLimit,
		originalLimit: c.originalLimit,
	}
}

type RandomChainGenerator struct {
	generators    []Generator
	limit         int
	originalLimit int
}

func (c RandomChainGenerator) Next() (interface{}, bool) {
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

func (c RandomChainGenerator) Iter() Generator {
	return RandomChainGenerator{
		generators:    c.generators,
		limit:         c.originalLimit,
		originalLimit: c.originalLimit,
	}
}

func ChainRoundRobin(limit int, g ...Generator) Generator {
	return RoundRobinChainGenerator{
		current:       0,
		generators:    g,
		limit:         limit,
		originalLimit: limit,
	}
}
func ChainRandom(limit int, g ...Generator) Generator {
	return RandomChainGenerator{
		generators:    g,
		limit:         limit,
		originalLimit: limit,
	}
}
