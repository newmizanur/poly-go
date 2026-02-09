パッケージ main

インポート (
    "fmt"
)

関数 main() {
    value := 5
    fmt.Println("value =", value)

    もし 真 && (value > 0) {
        fmt.Println("inside")
    } 違えば {
        fmt.Println("outside")
    }

    繰り返す i := 0; i < 3; i = i + 1 {
        もし i == 1 {
            次へ
        }
        fmt.Println("i", i)
    }

    条件分岐 value {
    場合 1:
        fmt.Println("one")
    場合 5:
        fmt.Println("five")
        抜ける
    その他:
        fmt.Println("other")
    }

    ch := 作る(チャネル 整数, 1)
    並行 send(ch)
    got := <-ch
    fmt.Println("got", got, "len=", 長さ([]整数{1,2,3}))
}

関数 send(ch チャネル 整数) {
    ch <- 42
}
