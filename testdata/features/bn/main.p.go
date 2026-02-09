প্যাকেজ main

আমদানি (
    "fmt"
)

ফাংশন main() {
    value := 5
    fmt.Println("value =", value)

    যদি সত্য && (value > 0) {
        fmt.Println("inside")
    } না_হলে {
        fmt.Println("outside")
    }

    জন্য i := 0; i < 3; i = i + 1 {
        যদি i == 1 {
            পরেরটায়_যাও
        }
        fmt.Println("i", i)
    }

    বাছাই value {
    বিকল্প 1:
        fmt.Println("one")
    বিকল্প 5:
        fmt.Println("five")
        থামো
    ডিফল্ট:
        fmt.Println("other")
    }

    ch := বানাও(চ্যানেল পূর্ণসংখ্যা, 1)
    চালাও send(ch)
    got := <-ch
    fmt.Println("got", got, "len=", দৈর্ঘ্য([]পূর্ণসংখ্যা{1,2,3}))
}

ফাংশন send(ch চ্যানেল পূর্ণসংখ্যা) {
    ch <- 42
}
