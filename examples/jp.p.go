//go:build linux || darwin
// +build linux darwin

パッケージ main

インポート (
	"fmt"
	"time"
)

定数 アプリ名 文字列 = "Japanese Poly-Go (Nadesiko-inspired)"
変数 グローバル整数 整数 = 0

型 ユーザー 構造体 {
	名前  文字列
	年齢 整数
}

型 あいさつ インタフェース {
	言う() 文字列
}

関数 (u ユーザー) 言う() 文字列 {
	戻す "こんにちは、" + u.名前
}

関数 main() {
	fmt.Println("🚀", アプリ名)

	// Strings/comments must not be translated:
	// もし, 違えば, 繰り返す, 反復, 条件分岐, 選択, 遅延, そのまま, 移動
	fmt.Println("string should remain untouched: もし, 違えば, 繰り返す, 条件分岐")

	// Escape prefix '@' demo (treat keyword-word as identifier)
	@型 := "これは識別子としての『型』"
	@構造体 := 123
	@インタフェース := 真
	fmt.Println("escaped:", 型, 構造体, インタフェース)

	// make/new/len/cap/append/copy/delete/close
	数 := 作る([]整数, 0, 8)
	fmt.Println("len/cap:", 長さ(数), 容量(数))

	数 = 追加(数, 1, 2, 3)

	宛先 := 作る([]整数, 3)
	複写数 := 複写(宛先, 数)
	fmt.Println("copied:", 複写数, 宛先)

	// map + delete
	m := 作る(辞書[文字列]整数)
	m["a"] = 1
	m["b"] = 2
	削除(m, "b")
	fmt.Println("map:", m)

	// new + interface usage
	u := 新規(ユーザー)
	u.名前 = "Rahim"
	u.年齢 = 30

	変数 g あいさつ = *u
	fmt.Println("greet:", g.言う())

	// for + range
	合計 := 0
	繰り返す _, v := 反復 数 {
		合計 += v
	}
	fmt.Println("sum:", 合計)

	// if/else + bool/nil/error
	変数 err 誤り = 無
	もし err == 無 && 真 && !偽 {
		fmt.Println("nil/bool ok")
	} 違えば {
		fmt.Println("unexpected")
	}

	// switch/case/default + fallthrough
	x := 1
	条件分岐 x {
	場合 1:
		fmt.Println("case 1")
		そのまま
	場合 2:
		fmt.Println("case 2 (via fallthrough)")
	その他:
		fmt.Println("default")
	}

	// goroutine + chan + select + defer
	ch := 作る(チャネル 文字列, 1)
	並行 仕事(ch)

	遅延 fmt.Println("defer executed at end of main")

	選択 {
	場合 msg := <-ch:
		fmt.Println("select recv:", msg)
	その他:
		fmt.Println("select default (no message yet)")
	}

	// break/continue
	繰り返す i := 0; i < 5; i++ {
		もし i == 2 {
			次へ
		}
		もし i == 4 {
			抜ける
		}
		fmt.Println("loop i:", i)
	}

	// goto demo
	もし 合計 > 0 {
		移動 終了
	}
	fmt.Println("this line should be skipped")

終了:
	fmt.Println("reached label: 終了")

	// panic/recover demo
	fmt.Println("panic/recover demo:", 安全呼び出し())

	// complex/real/imag demo
	z := 複素(2, 3)
	fmt.Println("complex:", z, "real:", 実部(z), "imag:", 虚部(z))

	fmt.Println("⏰ time:", time.Now())
	戻す
}

関数 仕事(out チャネル 文字列) {
	time.Sleep(50 * time.Millisecond)
	out <- "仕事完了"
	閉じる(out)
}

関数 安全呼び出し() 文字列 {
	遅延 関数() {
		_ = 回復()
	}()
	パニック("boom")
	戻す "unreachable"
}
