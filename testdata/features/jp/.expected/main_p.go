package main

import (
    "fmt"
)

func main() {
    value := 5
    fmt.Println("value =", value)

    if true && (value > 0) {
        fmt.Println("inside")
    } else {
        fmt.Println("outside")
    }

    for i := 0; i < 3; i = i + 1 {
        if i == 1 {
            continue
        }
        fmt.Println("i", i)
    }

    switch value {
    case 1:
        fmt.Println("one")
    case 5:
        fmt.Println("five")
        break
    default:
        fmt.Println("other")
    }

    ch := make(chan int, 1)
    go send(ch)
    got := <-ch
    fmt.Println("got", got, "len=", len([]int{1,2,3}))
}

func send(ch chan int) {
    ch <- 42
}
