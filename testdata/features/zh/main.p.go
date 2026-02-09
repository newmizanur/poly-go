包 main

导入 (
    "fmt"
)

函数 main() {
    value := 5
    fmt.Println("value =", value)

    如果 真 && (value > 0) {
        fmt.Println("inside")
    } 否则 {
        fmt.Println("outside")
    }

    循环 i := 0; i < 3; i = i + 1 {
        如果 i == 1 {
            继续
        }
        fmt.Println("i", i)
    }

    分支 value {
    情况 1:
        fmt.Println("one")
    情况 5:
        fmt.Println("five")
        跳出
    默认:
        fmt.Println("other")
    }

    ch := 创建(通道 整数, 1)
    并发 send(ch)
    got := <-ch
    fmt.Println("got", got, "len=", 长度([]整数{1,2,3}))
}

函数 send(ch 通道 整数) {
    ch <- 42
}
