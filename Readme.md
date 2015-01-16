# keen 

[![Build Status](https://travis-ci.org/gokeen/keen.svg?branch=master)][ci]
[![GoDoc](https://godoc.org/gopkg.in/gokeen/keen.v1?status.svg)][gd]

  [ci]: https://travis-ci.org/gokeen/keen
  [gd]: http://godoc.org/gopkg.in/gokeen/keen.v1

Keen.io in Go

## Install

    go get github.com/gokeen/keen

__gopkg.in__

    go get gopkg.in/gokeen/keen.v1

## Example

    k := keen.NewClient(func(c *keen.Client) {
      c.WriteKey = "awritekey"
    })

    err := k.Write(MyEvent{
      Action: "Wrote to Keen",
      Time:   time.Now(),
    })
    if err != nil {
      // handle error
    }

BYO-Event struct by implementing the `Event` interface.

    type MyEvent struct{
      Action string    `json:"action"`
      Time   time.Time `json:"time"`
    }

    func (MyEvent) ProjectID() string {
      return "aprojectid"
    }

    func (MyEvent) CollectionName() string {
      return "awesome-events"
    }


## License

MIT
