# fakenews

![fakenews.png](https://raw.githubusercontent.com/4thel00z/logos/master/fakenews.png)

## Motivation

I want to have a simple library which can generate me fake values in a fp manner.

## How to install

TBD

## Examples


### Infinite stream of fake data

```go
package main

import (
	"context"
	"github.com/4thel00z/fakenews/pkg/v1/fakenews"
	"log"
	"time"
)

type Message struct {
	Content string
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	in := make(chan Message)
	nameRule := fakenews.RandomPick[Message, string]("Content", "Hello", "World")
	message := Message{}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case val := <-in:
				println("Received:", val.Content)
			}
		}
	}()

	log.Fatalln(fakenews.InfiniteStream[Message](ctx, message, in, nameRule))
}
```

### Generate fake data based on a Reader

```go
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
```
### Single rule

```go
package main

import (
	"github.com/4thel00z/fakenews/pkg/v1/fakenews"
	"log"
)

type Person struct {
	Age  int
	Name string
}

func main() {
	// Create a rule to set a random integer between 1 and 10 for the Age field
	ageRule := fakenews.RandomInt[Person]("Age", 1, 10)
	nameRule := fakenews.RandomPick[Person, string]("Name", "John", "Jane", "Jack")
	fake, err := fakenews.Fake[Person](Person{}, ageRule, nameRule)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("My age is", fake.Age)
	log.Println("My name is", fake.Name)
}
```

## License

This project is licensed under the GPL-3 license.
