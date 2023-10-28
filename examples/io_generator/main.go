package main

import (
	"github.com/4thel00z/fakenews/pkg/v1/fakenews"
	"strings"
)

type Person struct {
}

func main() {
	loremIpsumGenerator := fakenews.GeneratorFromIO[string](strings.NewReader(strings.Repeat(`Lorem ipsum dolor sit amet, consectetur adipiscing elit.`, 1000)), func(line []byte) string {
		return string(line)
	})

}
