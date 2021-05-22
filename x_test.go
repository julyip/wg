//+build ignore

package x

import (
	. "github.com/ktye/wg/module"
	"github.com/ktye/wg/wasi_unstable"
)

// Init sets up the module. The function is not compiled to wasm.
func init() {
	Memory(1)
	Functions(0, Add, dup)
	Export(dup, ignore)
	Data(0, "abc")
}

// (func $Add (param $x i32) (param $y i32) (result i32)
// local.get $x local.get $y i32.add)
func Add(x, y int32) int32 { return x + y }

// (func $niladic)
func niladic() {}

// (func $dup (export "dup") (param $x i32) (result i32) (result i32)
// local.get $x local.get $x)
func dup(x int32) (int32, int32) { return x, x }

// (func $ignore (export "ignore") (param $x i32)
// local.get $x call $dup drop drop)
func ignore(x int32) { dup(x) }

// (func $statements call $niladic call $niladic)
func statements() { niladic(); niladic() }

// (func $rel (param $x i32) (param $y i32) (result i32)
// local.get $x i32.const 0 i32.gt_s local.get $y
// i32.const 0 i32.lt_s i32.and)
func rel(x, y int32) bool { return x > 0 && y < 0 }

// custom types
type I int64

// (func $negi (param $x i64) (result i64)
// i64.const 0 local.get $x i64.sub)
func negi(x I) I { return -x }

// (func $negf (param $x f32) (result f32)
// local.get $x f32.neg)
func negf(x float32) float32 { return -x }

// (func $clz (param $x i32) (result i32)
// local.get $x i32.clz)
func clz(x uint32) int32 { return I32clz(x) }

// structs
type st struct {
	a I
	b float64
}

// (func $structfn (param $s.a i64) (param $s.b f64)
// (param $y i64) (param $z f64) (result i64) (result f64)
// local.get $s.a local.get $y i64.add local.set $s.a
// local.get $s.b local.get $z f64.mul local.set $s.b
// local.get $s.a local.get $s.b)
func structfn(s st, y int64, z float64) st {
	s.a += I(y)
	s.b *= z
	return s
}

// (func $st.method (param $s.a i64) (param $s.b f64) (param $y i64) (result i64)
// local.get $s.a local.get $y i64.div_s)
func (s st) method(y int64) int64 { return int64(s.a) / y }

// (func $call (param $x f32) (result f32)
// local.get $x call $negf)
func call(x float32) float32 { return negf(x) }

// (func $call2 (param $x i32) (result i32) local.get $x call $dup call $Add)
func call2(x int32) int32 { return Add(dup(x)) }

// (func $call3 (param $x i32) (result i32) (local $a i32) (local $b i32)
// local.get $x call $dup local.set $b local.tee $a local.get $b i32.add)
func call3(x int32) int32 {
	a, b := dup(x)
	return a + b
}

// (func $retval (param $x i32) (result i32) (local $r i32)
// i32.const 1 local.get $x i32.add local.tee $r)
func retval(x int32) (r int32) {
	r = 1 + x
	return r
}

// (func $trap (param $x i32) (result i32)
// unreachable local.get $x)
func trap(x int32) int32 {
	panic(x)
	return x
}

// embed
type st1 struct {
	a int32
}
type st2 struct {
	st1
	a int64
}
type st3 struct {
	b st1
	a int64
}

// (func $st3.m (param $s.b.a i32) (param $s.a i64))
func (s st3) m() {}

// (func $st1.method1 (param $s.a i32) (result i32)
// i32.const 2 local.get $s.a i32.mul)
func (s st1) method1() int32 { return 2 * s.a }

// (func $callmethod1 (param $s.a i32) (result i32)
// local.get $s.a call $st1.method1)
func callmethod1(s st1) int32 { return s.method1() }

// (func $callinner1 (param $s.st1.a i32) (param $s.a i64) (result i32)
// local.get $s.st1.a call $st1.method1)
func callinner1(s st2) int32 { return s.st1.method1() }

/* not supported: direct method call of an embedded field
// (func $callinner2 (param $s.st1.a i32) (param $s.a i64) (result i32)
// local.get $s.st1.a call $st1.method1)
func callinner2(s st2) int32 { return s.method1() }
*/

type f2 func(int32, int32) int32

// (func $indirect (param $x i32) (param $y i32) (param $z i32) (result i32)
// local.get $y local.get $z local.get $x call_indirect (param i32) (param i32) (result i32))
func indirect(x, y, z int32) int32 {
	return Func[x].(f2)(y, z)
}

// (func $locals (param $x i32) (result i32) (local $a i32) (local $b i32)
// i32.const 0 local.get $x i32.sub local.set $b
// i32.const 2 local.get $b i32.mul local.tee $a)
func locals(x int32) int32 {
	var a int32
	b := -x
	a = 2 * b
	return a
}

// (func $varassign (param $x i32) (result i32) (local $y i32)
// i32.const 3 local.set $y local.get $x local.get $y i32.add)
func varassign(x int32) int32 {
	var y int32 = 3
	return x + y
}

// (func $varassign2 (param $x i32) (result i32) (local $a i32) (local $b i32)
// local.get $x i32.const 2 local.get $x i32.mul local.set $b local.tee $a local.get $b i32.sub)
func varassign2(x int32) int32 {
	var a, b int32 = x, 2 * x
	return a - b
}

// (func $localstruct (param $x i64) (result i32) (result i64)
// (local $s.st1.a i32) (local $s.a i64)
// local.get $x local.set $s.a local.get $s.st1.a local.get $s.a)
func localstruct(x int64) st2 {
	var s st2
	s.a = x
	return s
}

// (func $localstruct2 (param $x i64) (result i64)
// (local $s.st1.a i32) (local $s.a i64)
// local.get $x call $localstruct local.set $s.a local.set $s.st1.a local.get $s.a)
func localstruct2(x int64) int64 {
	var s st2 = localstruct(x)
	return s.a
}

var g int32
var g2 = int32(0)

// (func $globalasn (param $x i32) global.get $g local.get $x i32.add global.set $g)
func globalasn(x int32) { g += x }

// (func $globalvar (param $x i32) (result i32) local.get $x global.get $g i32.add)
func globalvar(x int32) int32 { return x + g }

var St3 st3
var Ge1, Ge2 = int64(1), int32(3)

// (func $globalstruct (param $x i32) (result i32) local.get $x global.get $St3.b.a i32.add)
func globalstruct(x int32) int32 { return x + St3.b.a }

const con = int32(3 << 5)

const (
	enum I = 1 << iota
	enum2
	enum3
)

// (func $constant (param $x i32) (result i32)
// local.get $x global.get $con i32.add)
func constant(x int32) int32 { return x + con }

// (func $enums (result i64)
// global.get $enum global.get $enum2 i64.add global.get $enum3 i64.add)
func enums() I { return enum + enum2 + enum3 }

// (func $load (param $addr i32) (result i32)
// local.get $addr i32.load local.get $addr i32.load8_s i32.add)
func load(addr int32) int32 { return I32(addr) + I8(addr) }

// (func $store (param $addr i32) (param $x i32)
// local.get $addr local.get $x i32.store)
func store(addr int32, x int32) {
	SetI32(addr, x)
}

// (func $iff (param $x i32) (result i32)
// local.get $x i32.const 0 i32.gt_s if i32.const 2 local.get $x i32.mul return end local.get $x)
func iff(x int32) int32 {
	if x > 0 {
		return 2 * x
	}
	return x
}

// (func $ifelse (param $x i32) (result i32)
// local.get $x i32.const 0 i32.gt_s if i32.const 2 local.get $x i32.mul return
// else local.get $x return end local.get $x)
func ifelse(x int32) int32 {
	if x > 0 {
		return 2 * x
	} else {
		return x
	}
	return x // unreached
}

// (func $elseif (param $x i32) (result i32)
// local.get $x i32.const 0 i32.gt_s if local.get $x i32.const 1 i32.add local.set $x
// else local.get $x i32.const 0 i32.lt_s if local.get $x i32.const 1 i32.sub local.set $x end end local.get $x)
func elseif(x int32) int32 {
	if x > 0 {
		x++
	} else if x < 0 {
		x--
	}
	return x
}

// (func $ifinit (param $x i32) (result i32) (local $n.1 i32)
// local.get $x local.get $x i32.mul local.tee $n.1 i32.const 0 i32.gt_s
// if i32.const 0 local.get $n.1 i32.sub return end local.get $x)
func ifinit(x int32) int32 {
	if n := x * x; n > 0 {
		return -n
	}
	return x
}

// (func $while (param $n i32) (result i32) (local $r i32)
// block loop local.get $n i32.const 0 i32.gt_s i32.eqz br_if 1
// local.get $r i32.const 1 i32.add local.set $r br 0 end end
// local.get $r)
func while(n int32) (r int32) {
	for n > 0 {
		r++
	}
	return r
}

// (func $forloop (param $n i32) (result i32) (local $r i32) (local $i.1 i32)
// i32.const 0 local.set $i.1 block loop local.get $i.1 local.get $n i32.lt_s i32.eqz br_if 1
// local.get $r local.get $i.1 i32.add local.set $r local.get $i.1 i32.const 1 i32.add local.set $i.1 br 0 end end local.get $r)
func forloop(n int32) (r int32) {
	for i := int32(0); i < n; i++ {
		r += i
	}
	return r
}

// (func $forbreak (param $n i32) (result i32) (local $r i32)
// block loop local.get $r local.get $n i32.eq if br 2 end
// local.get $r i32.const 1 i32.add local.set $r br 0 end end local.get $r)
func forbreak(n int32) (r int32) {
	for {
		if r == n {
			break
		}
		r++
	}
	return r
}

// (func $forcontinue (param $n i32) (result i32) (local $r i32)
// block loop
// local.get $r i32.const 5 i32.eq if br 1 end
// local.get $r local.get $n i32.gt_s if br 2 end
// local.get $r i32.const 1 i32.add local.set $r br 0 end end local.get $r)
func forcontinue(n int32) (r int32) {
	for {
		if r == 5 {
			continue
		}
		if r > n {
			break
		}
		r++
	}
	return r
}

// (func $forlabel (param $n i32) (result i32) (local $r i32) (local $i.1 i32)
// i32.const 0 local.set $i.1
// block $out:1 loop $out:2 local.get $i.1 local.get $n i32.lt_s i32.eqz br_if 1
// block loop local.get $r i32.const 5 i32.eq if br $out:1 end
// local.get $r   i32.const 1 i32.add local.set $r br 0 end end
// local.get $i.1 i32.const 1 i32.add local.set $i.1 br 0 end end local.get $r)
func forlabel(n int32) (r int32) {
out:
	for i := int32(0); i < n; i++ {
		for {
			if r == 5 {
				break out
			}
			r++
		}
	}
	return r
}

// (func $brtable (param $x i32) (result i32)
// block block block
// i32.const 1 local.get $x i32.add
// br_table 0 1 2 br 2 end
// local.get $x i32.const 2 i32.add local.set $x br 1 end
// local.get $x i32.const 2 i32.mul local.set $x br 0 end local.get $x)
func brtable(x int32) int32 {
	switch uint32(1 + x) {
	case 0:
	case 1:
		x += 2
	default:
		x *= 2
	}
	return x
}

// (func $scope (result i32) (local $i i32) (local $i.1 i32)
// i32.const 1 local.set $i
// i32.const 0 local.set $i.1
// block loop local.get $i.1 i32.const 5 i32.lt_s i32.eqz br_if 1
// local.get $i.1 i32.const 1 i32.add local.set $i.1 br 0 end end local.get $i)
func scope() int32 {
	i := int32(1)
	for i := int32(0); i < 5; i++ {
	}
	return i
}

// (func $scopevar (param $x i32) (result i32) (local $i.2 i32)
// local.get $x i32.const 0 i32.gt_s if i32.const 0 local.set $i.2
// local.get $x i32.const 0 i32.lt_s if local.get $i.2 local.set $x end end local.get $x)
func scopevar(x int32) int32 {
	if x > 0 {
		var i int32 = 0
		if x < 0 {
			x = i
		}
	}
	return x
}

// (func $memcpy (param $x i32)
// local.get $x local.get $x i32.const 10 i32.add i32.const 5 memory.copy)
func memcpy(x int32) {
	Memorycopy(x, x+10, 5)
}

// (func $wasicall (result i64)
// i32.const 0 i64.const 0 call $wasi_unstable.clock_time_get)
func wasicall() wasi_unstable.Timestamp { return wasi_unstable.Clock_time_get(0, 0) }
