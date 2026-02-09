paquete main

importar (
    "fmt"
)

funcion main() {
    value := 5
    fmt.Println("value =", value)

    si verdadero && (value > 0) {
        fmt.Println("inside")
    } sino {
        fmt.Println("outside")
    }

    para i := 0; i < 3; i = i + 1 {
        si i == 1 {
            siguiente
        }
        fmt.Println("i", i)
    }

    cambiar value {
    caso 1:
        fmt.Println("one")
    caso 5:
        fmt.Println("five")
        romper
    defecto:
        fmt.Println("other")
    }

    ch := crear(canal entero, 1)
    ir send(ch)
    got := <-ch
    fmt.Println("got", got, "len=", longitud([]entero{1,2,3}))
}

funcion send(ch canal entero) {
    ch <- 42
}
