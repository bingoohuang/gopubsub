# 🚌 pubsub

[![Build Status](https://travis-ci.org/bingoohuang/gopubsub.svg?branch=master)](https://travis-ci.org/bingoohuang/gopubsub)
[![Go Report Card](https://goreportcard.com/badge/github.com/bingoohuang/gopubsub)](https://goreportcard.com/report/github.com/bingoohuang/gopubsub)
[![codecov](https://codecov.io/gh/bingoohuang/gopubsub/branch/master/graph/badge.svg)](https://codecov.io/gh/bingoohuang/gopubsub)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fvardius%2Fmessage-bus.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fvardius%2Fmessage-bus?ref=badge_shield)
[![](https://godoc.org/github.com/bingoohuang/gopubsub?status.svg)](http://godoc.org/github.com/bingoohuang/gopubsub)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/bingoohuang/gopubsub/blob/master/LICENSE.md)

Go simple async message bus.


## 🚅 Benchmark

Time complexity of a `Pub` method is considered to be [linear time `O(n)`](https://en.wikipedia.org/wiki/Time_complexity#Linear_time). Where **n** corresponds to the number of *subscribers* for a given **topic**.

```bash
➜  message-bus git:(master) ✗  system_profiler SPHardwareDataType Hardware:
Hardware:

    Hardware Overview:

      Model Name: MacBook Pro
      Model Identifier: MacBookPro15,1
      Processor Name: 6-Core Intel Core i7
      Processor Speed: 2.2 GHz
      Number of Processors: 1
      Total Number of Cores: 6
      L2 Cache (per Core): 256 KB
      L3 Cache: 9 MB
      Hyper-Threading Technology: Enabled
      Memory: 16 GB
      Boot ROM Version: 1037.80.53.0.0 (iBridge: 17.16.13050.0.0,0)
      Activation Lock Status: Disabled
```

```bash
➜  gopubsub git:(master) ✗ go test -bench=. -cpu=4 -benchmem
goos: darwin
goarch: amd64
pkg: github.com/bingoohuang/gopubsub
BenchmarkPublish-4                       4843279               251 ns/op               0 B/op          0 allocs/op
BenchmarkSubscribe-4                      826225              1824 ns/op             768 B/op          5 allocs/op

```

## 🏫 Basic example
```go
package main

import (
    "fmt"
    "sync"

    "github.com/bingoohuang/gopubsub"
)

func main() {
    queueSize := 100
    bus := gopubsub.New(queueSize)

    var wg sync.WaitGroup
    wg.Add(2)

    _ = bus.Sub("topic", func(v bool) {
        defer wg.Done()
        fmt.Println(v)
    })

    _ = bus.Sub("topic", func(v bool) {
        defer wg.Done()
        fmt.Println(v)
    })

    // Publish block only when the buffer of one of the subscribers is full.
    // change the buffer size altering queueSize when creating new messagebus
    bus.Pub("topic", true)
    wg.Wait()
}
```
