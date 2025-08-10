package main

import (
	"sync"
	"time"
	"fmt"
)

func or(channels ...<-chan interface{}) <-chan interface{} {    
    switch len(channels) {
    case 0:
        return nil
    case 1:
        return channels[0]
    }

    orDone := make(chan interface{})
        
    go func() {        
        var once sync.Once
                
        for _, ch := range channels {
            go func(c <-chan interface{}) {
                select {
                case <-c:                    
                    once.Do(func() { close(orDone) })
                case <-orDone:                    
                }
            }(ch)
        }
                
        <-orDone
    }()
    
    return orDone
}

func sig(after time.Duration) <-chan interface{} {
    c := make(chan interface{})
    go func() {
        defer close(c)
        time.Sleep(after)
    }()
    return c
}

func main() {
    start := time.Now()
    
    <-or(
        sig(2*time.Hour),
        sig(5*time.Minute),
        sig(1*time.Second),
        sig(1*time.Hour),
        sig(1*time.Minute),
    )
    
    fmt.Printf("done after %v\n", time.Since(start))
}