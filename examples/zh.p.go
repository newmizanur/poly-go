//go:build linux || darwin
// +build linux darwin

包 main

导入 (
	"fmt"
	"time"
)

常量 应用名 字符串 = "Chinese Poly-Go (Wenyan-inspired)"
变量 全局计数 整数 = 0

类型 用户 结构体 {
	名字 字符串
	年龄 整数
}

类型 问候 接口 {
	说() 字符串
}

函数 (u 用户) 说() 字符串 {
	返回 "你好，" + u.名字
}

函数 main() {
	fmt.Println("🚀", 应用名)

	// Strings/comments must not be translated:
	// 如果, 否则, 循环, 遍历, 分支, 选择, 延迟, 贯穿, 跳转
	fmt.Println("string should remain untouched: 如果, 否则, 循环, 分支")

	// Escape prefix '@' demo (treat keyword-word as identifier)
	@类型 := "这是名为『类型』的普通变量"
	@结构体 := 123
	@接口 := 真
	fmt.Println("escaped:", 类型, 结构体, 接口)

	// make/new/len/cap/append/copy/delete/close
	数字 := 创建([]整数, 0, 8)
	fmt.Println("len/cap:", 长度(数字), 容量(数字))

	数字 = 追加(数字, 1, 2, 3)

	目标 := 创建([]整数, 3)
	已复制 := 复制(目标, 数字)
	fmt.Println("copied:", 已复制, 目标)

	// map + delete
	m := 创建(映射[字符串]整数)
	m["a"] = 1
	m["b"] = 2
	删除(m, "b")
	fmt.Println("map:", m)

	// new + interface usage
	u := 新建(用户)
	u.名字 = "Rahim"
	u.年龄 = 30

	变量 g 问候 = *u
	fmt.Println("greet:", g.说())

	// for + range
	总和 := 0
	循环 _, v := 遍历 数字 {
		总和 += v
	}
	fmt.Println("sum:", 总和)

	// if/else + bool/nil/error
	变量 err 错误 = 空
	如果 err == 空 && 真 && !假 {
		fmt.Println("nil/bool ok")
	} 否则 {
		fmt.Println("unexpected")
	}

	// switch/case/default + fallthrough
	x := 1
	分支 x {
	情况 1:
		fmt.Println("case 1")
		贯穿
	情况 2:
		fmt.Println("case 2 (via fallthrough)")
	默认:
		fmt.Println("default")
	}

	// goroutine + chan + select + defer
	ch := 创建(通道 字符串, 1)
	并发 工作(ch)

	延迟 fmt.Println("defer executed at end of main")

	选择 {
	情况 msg := <-ch:
		fmt.Println("select recv:", msg)
	默认:
		fmt.Println("select default (no message yet)")
	}

	// break/continue
	循环 i := 0; i < 5; i++ {
		如果 i == 2 {
			继续
		}
		如果 i == 4 {
			跳出
		}
		fmt.Println("loop i:", i)
	}

	// goto demo
	如果 总和 > 0 {
		跳转 结束
	}
	fmt.Println("this line should be skipped")

结束:
	fmt.Println("reached label: 结束")

	// panic/recover demo
	fmt.Println("panic/recover demo:", 安全调用())

	// complex/real/imag demo
	z := 复数(2, 3)
	fmt.Println("complex:", z, "real:", 实部(z), "imag:", 虚部(z))

	fmt.Println("⏰ time:", time.Now())
	返回
}

函数 工作(out 通道 字符串) {
	time.Sleep(50 * time.Millisecond)
	out <- "完成"
	关闭(out)
}

函数 安全调用() 字符串 {
	延迟 函数() {
		_ = 恢复()
	}()
	恐慌("boom")
	返回 "unreachable"
}
