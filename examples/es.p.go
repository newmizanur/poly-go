//go:build linux || darwin
// +build linux darwin

paquete main

importar (
	"fmt"
	"time"
)

// ====== Global const/var ======
constante app_nombre cadena = "Spanish-Go Full Coverage"
var contador_global entero = 0

// ====== tipo/estructura/interfaz ======
tipo usuario estructura {
	nombre cadena
	edad   entero
}

tipo saludador interfaz {
	saludar() cadena
}

funcion (u usuario) saludar() cadena {
	retornar "Hola, " + u.nombre
}

// ====== main ======
funcion main() {
	fmt.Println("🚀", app_nombre)

	// Strings/comments must not be translated:
	// si, para, retornar, tipo, estructura, interfaz, rango, aplazar, seleccionar, continuar_caso, ir_a
	fmt.Println("string should remain untouched: si, para, retornar, tipo, estructura, interfaz")

	// ====== Escape prefix '@' ======
	@tipo := "esto es un identificador llamado 'tipo'"
	@estructura := 123
	@interfaz := verdadero
	fmt.Println("escaped:", tipo, estructura, interfaz)

	// ====== make/new/len/cap/append/copy/delete/close ======
	numeros := crear([]entero, 0, 8)
	fmt.Println("len/cap:", longitud(numeros), capacidad(numeros))

	numeros = agregar(numeros, 1, 2, 3)

	destino := crear([]entero, 3)
	copiados := copiar(destino, numeros) // copy(dst, src)
	fmt.Println("copied:", copiados, destino)

	// map keyword + delete builtin
	m := crear(mapa[cadena]entero)
	m["a"] = 1
	m["b"] = 2
	borrar(m, "b")
	fmt.Println("map:", m)

	// new builtin
	u := nuevo(usuario)
	u.nombre = "Rahim"
	u.edad = 30

	// interface usage
	var g saludador = *u
	fmt.Println("greet:", g.saludar())

	// ====== for + range ======
	suma := 0
	para _, v := rango numeros {
		suma += v
	}
	fmt.Println("sum:", suma)

	// ====== if/else + bool/nil/error ======
	var err error = nulo
	si err == nulo && verdadero && !falso {
		fmt.Println("nil/bool ok")
	} sino {
		fmt.Println("unexpected")
	}

	// ====== switch/case/default + fallthrough ======
	x := 1
	cambiar x {
	caso 1:
		fmt.Println("case 1")
		continuar_caso
	caso 2:
		fmt.Println("case 2 (via fallthrough)")
	defecto:
		fmt.Println("default")
	}

	// ====== goroutine + chan + select + defer ======
	ch := crear(canal cadena, 1)
	ir trabajo(ch)

	aplazar fmt.Println("defer executed at end of main")

	seleccionar {
	caso msg := <-ch:
		fmt.Println("select recv:", msg)
	defecto:
		fmt.Println("select default (no message yet)")
	}

	// ====== break/continue ======
	para i := 0; i < 5; i++ {
		si i == 2 {
			siguiente
		}
		si i == 4 {
			romper
		}
		fmt.Println("loop i:", i)
	}

	// ====== goto (keyword) demo ======
	si suma > 0 {
		ir_a fin_evento
	}
	fmt.Println("this line should be skipped by goto")

fin_evento:
	fmt.Println("reached label: fin_evento")

	// ====== panic/recover (predeclared) demo ======
	fmt.Println("panic/recover demo:")
	fmt.Println("safeCall result:", llamada_segura())

	// ====== complex/real/imag demo (predeclared) ======
	z := complejo(2, 3) // complex(2,3)
	fmt.Println("complex:", z, "real:", real(z), "imag:", imaginario(z))

	fmt.Println("⏰ time:", time.Now())
	retornar
}

funcion trabajo(out canal cadena) {
	time.Sleep(50 * time.Millisecond)
	out <- "trabajo terminado"
	cerrar(out)
}

// ====== panic/recover helpers ======
funcion llamada_segura() cadena {
	aplazar funcion() {
		_ = recuperar() // recover()
	}()
	panico("boom") // panic("boom")
	retornar "unreachable"
}
