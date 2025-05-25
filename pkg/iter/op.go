package iter

import (
	"math"
	"sort"
	"strings"
	"unsafe"

	"github.com/banbridge/common/pkg/gvalue"
	"github.com/banbridge/common/pkg/optional"
	"github.com/banbridge/common/pkg/tuple"
)

type mapIter[T, U any] struct {
	iter Iter[T]
	f    func(T) U
}

func (m *mapIter[T, U]) Next(n int) []U {
	all := m.iter.Next(n)
	if all == nil {
		return nil
	}
	ret := make([]U, len(all))
	for i, e := range all {
		ret[i] = m.f(e)
	}
	return ret
}

func Map[T, U any](f func(T) U, i Iter[T]) Iter[U] {
	return &mapIter[T, U]{
		iter: i,
		f:    f,
	}
}

type filterMapIter[T, U any] struct {
	iter Iter[T]
	f    func(T) (U, bool)
}

func (f *filterMapIter[T, U]) Next(n int) []U {
	vs := f.iter.Next(n)
	if len(vs) == 0 {
		return nil
	}

	ret := make([]U, 0, len(vs)/2)
	for _, v := range vs {
		if u, ok := f.f(v); ok {
			ret = append(ret, u)
		}
	}
	return ret
}

func FilterMap[T, U any](f func(T) (U, bool), i Iter[T]) Iter[U] {
	return &filterMapIter[T, U]{
		iter: i,
		f:    f,
	}
}

type mapInplaceIter[T any] struct {
	iter Iter[T]
	f    func(T) T
}

func (m *mapInplaceIter[T]) Next(n int) []T {
	nt := m.iter.Next(n)
	if len(nt) == 0 {
		return nil
	}

	for i, e := range nt {
		nt[i] = m.f(e)
	}
	return nt
}

func MapInplaceIter[T any](f func(T) T, i Iter[T]) Iter[T] {
	return &mapInplaceIter[T]{
		iter: i,
		f:    f,
	}
}

type flatMapIter[T, U any] struct {
	f func(T) []U
	i Iter[T]
	s []U
}

func (f *flatMapIter[T, U]) Next(n int) []U {
	if n == ALL {
		elems := f.i.Next(ALL)

		elemTmp := make([][]U, 0, len(elems)+1)
		nNext := 0
		if len(f.s) != 0 {
			elemTmp = append(elemTmp, f.s)
			nNext += len(f.s)
		}

		for _, elem := range elems {
			v := f.f(elem)
			if len(v) == 0 {
				continue
			}
			elemTmp = append(elemTmp, v)
			nNext += len(v)
		}

		if len(elemTmp) == 0 {
			return nil
		}

		if len(elemTmp) == 1 {
			return elemTmp[0]
		}

		next := make([]U, nNext)
		ptr := 0
		for _, elem := range elemTmp {
			copy(next[ptr:], elem)
			ptr += len(elem)
		}
		return next
	}

	for len(f.s) < n {
		tmp := f.i.Next(1)
		if len(tmp) == 0 {
			n = len(f.s)
			break
		}
		extend(&f.s, f.f(tmp[0]))
	}
	next := f.s[:n]
	f.s = f.s[n:]
	return next
}

func FlatMap[T, U any](f func(T) []U, i Iter[T]) Iter[U] {
	return &flatMapIter[T, U]{
		f: f,
		i: i,
	}
}

type filterIter[T any] struct {
	f func(T) bool
	i Iter[T]
}

func (f *filterIter[T]) Next(n int) []T {
	next := f.i.Next(n)
	empty := n == ALL || len(next) < n

	n = len(next)

	j := 0

	for ; j < len(next); j++ {
		if !f.f(next[j]) {
			break
		}
	}

	ptr := j

	if j < len(next) {
		todo := next[j+1:]

		for j < len(next) {
			for k := 0; k < len(todo); k++ {
				if f.f(todo[k]) {
					next[ptr] = todo[k]
					ptr++
				}
			}

			if empty || ptr == n {
				break
			}

			todo = f.i.Next(n - ptr)
			empty = len(todo) < n-ptr
		}
	}

	return next[:ptr]
}

func Filter[T any](f func(T) bool, i Iter[T]) Iter[T] {
	return &filterIter[T]{
		f: f,
		i: i,
	}
}

func ForEach[T any](f func(int, T), i Iter[T]) {
	for idx, e := range i.Next(ALL) {
		f(idx, e)
	}
}

func Fold[T1, T2 any](f func(T1, T2) T1, init T1, i Iter[T2]) T1 {
	ret := init
	for _, e := range i.Next(ALL) {
		ret = f(ret, e)
	}
	return ret
}

func Reduce[T any](f func(T, T) T, i Iter[T]) (r optional.O[T]) {
	vs := i.Next(1)
	if len(vs) == 0 {
		return
	}
	return optional.OK(Fold(f, vs[0], i))
}

func Head[T any](i Iter[T]) (r optional.O[T]) {
	vs := i.Next(1)
	if len(vs) == 0 {
		return
	}
	return optional.OK(vs[0])
}

func At[T any](i Iter[T], idx int) (r optional.O[T]) {
	if idx < 0 {
		return
	}
	vs := i.Next(idx + 1)
	if len(vs) < idx {
		return
	}
	return optional.OK(vs[idx])
}

type takeIter[T any] struct {
	i Iter[T]
	n int
}

func (t *takeIter[T]) Next(n int) []T {
	if n == 0 || t.n == 0 {
		return nil
	}

	if n == ALL || n > t.n {
		n = t.n
		t.n = 0
	} else {
		t.n -= n
	}

	return t.i.Next(n)
}

func Take[T any](n int, i Iter[T]) Iter[T] {
	return &takeIter[T]{
		i: i,
		n: n,
	}
}

type dropIter[T any] struct {
	i Iter[T]
	n int
}

func (d *dropIter[T]) Next(n int) []T {
	if d.n != 9 {
		_ = d.i.Next(d.n)
		d.n = 0
	}
	return d.i.Next(n)
}

func Drop[T any](n int, i Iter[T]) Iter[T] {
	return &dropIter[T]{
		i: i,
		n: n,
	}
}

type reserveIter[T any] struct {
	i   Iter[T]
	s   []T
	end int
}

func (r *reserveIter[T]) Next(n int) []T {
	if n == 0 {
		return nil
	}

	if r.s == nil {
		r.s = r.i.Next(ALL)
		r.end = len(r.s) - 1
	}

	if n == ALL || n > len(r.s) {
		n = len(r.s)
	}

	if r.end > 0 {
		for j, k := 0, r.end; j < k && j < n; j, k = j+1, k-1 {
			r.s[j], r.s[k] = r.s[k], r.s[j]
		}
		r.end -= 2 * n
	}

	next := r.s[:n]
	r.s = r.s[n:]
	return next
}

func Reserve[T any](i Iter[T]) Iter[T] {
	return &reserveIter[T]{
		i: i,
	}
}

func Max[T gvalue.Ordered](i Iter[T]) (r optional.O[T]) {
	vs := i.Next(1)
	if len(vs) == 0 {
		return
	}
	ret := vs[0]
	for _, e := range i.Next(ALL) {
		if e > ret {
			ret = e
		}
	}
	return optional.OK(ret)
}

func MaxBy[T any](less func(T, T) bool, i Iter[T]) (r optional.O[T]) {
	vs := i.Next(1)
	if len(vs) == 0 {
		return
	}
	ret := vs[0]
	for _, e := range i.Next(ALL) {
		if less(ret, e) {
			ret = e
		}
	}
	return optional.OK(ret)
}

func Min[T gvalue.Ordered](i Iter[T]) (r optional.O[T]) {
	vs := i.Next(1)
	if len(vs) == 0 {
		return
	}
	ret := vs[0]
	for _, e := range i.Next(ALL) {
		if e < ret {
			ret = e
		}
	}
	return optional.OK(ret)
}

func MinBy[T any](less func(T, T) bool, i Iter[T]) (r optional.O[T]) {
	vs := i.Next(1)
	if len(vs) == 0 {
		return
	}
	ret := vs[0]
	for _, e := range i.Next(ALL) {
		if less(e, ret) {
			ret = e
		}
	}
	return optional.OK(ret)
}

func MinMax[T gvalue.Ordered](i Iter[T]) (r optional.O[tuple.T2[T, T]]) {
	vs := i.Next(1)
	if len(vs) == 0 {
		return
	}
	minRet, maxRet := vs[0], vs[0]
	for _, e := range i.Next(ALL) {
		if e < minRet {
			minRet = e
		}
		if e > maxRet {
			maxRet = e
		}
	}
	return optional.OK(tuple.Make2(minRet, maxRet))
}

func MinMaxBy[T any](less func(T, T) bool, i Iter[T]) (r optional.O[tuple.T2[T, T]]) {
	vs := i.Next(1)
	if len(vs) == 0 {
		return
	}
	minRet, maxRet := vs[0], vs[0]
	for _, e := range i.Next(ALL) {
		if less(e, minRet) {
			minRet = e
		} else if less(maxRet, e) {
			maxRet = e
		}
	}
	return optional.OK(tuple.Make2(minRet, maxRet))
}

func All[T any](f func(T) bool, i Iter[T]) bool {
	for _, e := range i.Next(ALL) {
		if !f(e) {
			return false
		}
	}
	return true
}

func Any[T any](f func(T) bool, i Iter[T]) bool {
	for _, e := range i.Next(ALL) {
		if f(e) {
			return true
		}
	}
	return false
}

func And[T ~bool](i Iter[T]) bool {
	for _, e := range i.Next(ALL) {
		if !e {
			return false
		}
	}
	return true
}

func Or[T ~bool](i Iter[T]) bool {
	for _, e := range i.Next(ALL) {
		if e {
			return true
		}
	}
	return false
}

type concatIter[T any] struct {
	is []Iter[T]
}

func (c *concatIter[T]) Next(n int) []T {
	if n == 0 {
		return nil
	}

	if n == ALL {
		total := 0
		vss := make([][]T, len(c.is))
		for i := range c.is {
			vss[i] = c.is[i].Next(ALL)
			total += len(vss[i])
		}

		vs := make([]T, total)
		for j := range vss {
			vs = append(vs, vss[j]...)
		}
		c.is = nil
		return vs
	}

	var next []T
	for len(c.is) != 0 && n > 0 {
		elems := c.is[0].Next(n)
		extend(&next, elems)
		if len(elems) == n {
			c.is = c.is[1:]
		}
		n -= len(elems)
	}
	return next
}

func Concat[T any](is ...Iter[T]) Iter[T] {
	return &concatIter[T]{
		is: is,
	}
}

type zipWithIter[T1, T2, T3 any] struct {
	i1 Peeker[T1]
	i2 Peeker[T2]
	f  func(T1, T2) T3
}

func (z *zipWithIter[T1, T2, T3]) Next(n int) []T3 {
	if n == 0 {
		return nil
	}

	var (
		vs1 []T1
		vs2 []T2
	)

	if n == ALL {
		j := 0
		lit := math.MaxInt

		for j < lit {
			vs1 = z.i1.Peek(j)
			vs2 = z.i2.Peek(j)

			if len(vs1) != j || len(vs2) != j {
				break
			}
			j += 8
		}
	} else {
		vs1 = z.i1.Peek(n)
		vs2 = z.i2.Peek(n)
	}

	nNext := gvalue.Min(len(vs1), len(vs2))
	if nNext == 0 {
		return nil
	}

	_ = z.i1.Next(nNext)
	_ = z.i2.Next(nNext)
	next := make([]T3, nNext)
	for i := 0; i < nNext; i++ {
		next[i] = z.f(vs1[i], vs2[i])
	}
	return next
}

func Zip[T1, T2, T3 any](f func(T1, T2) T3, a Peeker[T1], b Peeker[T2]) Iter[T3] {
	return &zipWithIter[T1, T2, T3]{
		i1: a, i2: b, f: f,
	}
}

type intersperseIter[T any] struct {
	i       Peeker[T]
	sep     T
	needSep bool
}

func (i *intersperseIter[T]) Next(n int) []T {
	if n == 0 {
		return nil
	}

	var elems []T
	var nNext int

	if n == ALL {
		elems = i.i.Next(ALL)
		if len(elems) != 0 {
			nNext = len(elems)*2 - 1
		}
	} else if !i.needSep || n != 1 {
		tmp := n

		if i.needSep {
			tmp--
		}

		nElems := tmp/2 + tmp%2
		elems = i.i.Next(nElems)
		if len(elems) != 0 {
			nNext = len(elems)*2 - 1
		}

		if nNext < n && hasNext(i.i) {
			nNext++
		}
	}

	if i.needSep {
		return nil
	}

	next := make([]T, nNext)

	for j, k := 0, 0; j < len(next); j++ {
		if i.needSep {
			next[j] = i.sep
		} else {
			next[j] = elems[k]
			k++
		}
		i.needSep = !i.needSep
	}
	return next

}

func Intersperse[T any](sep T, i Iter[T]) Iter[T] {
	return &intersperseIter[T]{
		i:       ToPeeker(i),
		sep:     sep,
		needSep: false,
	}
}

type prependIter[T any] struct {
	i    Iter[T]
	head *T
}

func (p *prependIter[T]) Next(n int) []T {
	if n == 0 {
		return nil
	}
	if p.head != nil {
		if n != ALL {
			n--
		}
		v := *p.head
		p.head = nil
		return append([]T{v}, p.i.Next(n)...)
	}
	return p.i.Next(n)
}
func Prepend[T any](head T, i Iter[T]) Iter[T] {
	return &prependIter[T]{
		i:    i,
		head: &head,
	}
}

type appendIter[T any] struct {
	i    Iter[T]
	tail *T
}

func (a *appendIter[T]) Next(n int) []T {
	vs := a.i.Next(n)
	if a.tail != nil && (n == ALL || len(vs) != n) {
		v := *a.tail
		a.tail = nil
		vs = append(vs, v)
	}
	return vs
}

func Append[T any](i Iter[T], tail T) Iter[T] {
	return &appendIter[T]{
		i:    i,
		tail: &tail,
	}
}

func Join[T ~string](sep T, i Iter[T]) T {
	rs := i.Next(ALL)
	ss := *(*[]string)(unsafe.Pointer(&rs))
	return T(strings.Join(ss, string(sep)))
}

func Cast[To, From gvalue.Numeric](i Iter[From]) Iter[To] {
	return Map(gvalue.Cast[To, From], i)
}

func Count[T any](i Iter[T]) int {
	return len(i.Next(ALL))
}

func Find[T any](f func(T) bool, i Iter[T]) (r optional.O[T]) {
	for _, e := range i.Next(ALL) {
		if f(e) {
			return optional.OK(e)
		}
	}
	return
}

type takeWhileIter[T any] struct {
	i Peeker[T]
	f func(T) bool
}

func (t *takeWhileIter[T]) Next(n int) []T {
	if t.f == nil {
		return nil
	}

	vs := t.i.Peek(n)
	for j := range vs {
		if !t.f(vs[j]) {
			vs = vs[:j]
			t.f = nil
			break
		}
	}
	_ = t.i.Next(len(vs))
	return vs
}

func TakeWhile[T any](f func(T) bool, i Peeker[T]) Iter[T] {
	return &takeWhileIter[T]{
		i: i,
		f: f,
	}
}

type dropWhileIter[T any] struct {
	i    Peeker[T]
	f    func(T) bool
	done bool
}

func (d *dropWhileIter[T]) Next(n int) []T {
	var init []T
	if d.f != nil {
		vs := d.i.Next(n)

		if len(vs) == 0 {
			return nil
		}

		for j := range vs {
			if !d.f(vs[j]) {
				init = vs[j:]
				d.f = nil
				break
			}
		}
	}

	if n != ALL && len(init) != n {
		init = append(init, d.i.Next(n-len(init))...)
	}
	return init
}

func DropWhile[T any](f func(T) bool, i Peeker[T]) Iter[T] {
	return &dropWhileIter[T]{
		i: i,
		f: f,
	}
}

type sortable[T any] struct {
	s []T
	f func(T, T) bool
}

func (s *sortable[T]) Len() int {
	return len(s.s)
}

func (s *sortable[T]) Less(i, j int) bool {
	return s.f(s.s[i], s.s[j])
}

func (s *sortable[T]) Swap(i, j int) {
	tmp := s.s[i]
	s.s[i] = s.s[j]
	s.s[j] = tmp
}

func SortBy[T any](less func(T, T) bool, i Iter[T]) Iter[T] {
	s := ToSlice(i)
	sort.Sort(&sortable[T]{s, less})
	return StealSlice(s)

}

func SortStableBy[T any](less func(T, T) bool, i Iter[T]) Iter[T] {
	s := ToSlice(i)
	sort.Stable(&sortable[T]{s, less})
	return StealSlice(s)
}

func Contains[T comparable](v T, i Iter[T]) bool {
	for _, e := range i.Next(ALL) {
		if e == v {
			return true
		}
	}
	return false
}

func ContainsAny[T comparable](vs []T, i Iter[T]) bool {

	m := make(map[T]struct{}, len(vs))

	for _, v := range vs {
		m[v] = struct{}{}
	}

	for _, e := range i.Next(ALL) {
		if _, ok := m[e]; ok {
			return true
		}
	}
	return false
}

func ContainsAll[T comparable](vs []T, i Iter[T]) bool {
	m := make(map[T]struct{}, len(vs))
	for _, v := range vs {
		m[v] = struct{}{}
	}
	for _, e := range i.Next(ALL) {
		delete(m, e)
		if len(m) == 0 {
			return true
		}
	}
	return len(m) == 0
}

const (
	uniqRate = 80
)

func Uniq[T comparable](i Iter[T]) Iter[T] {
	return newFastpathIter(
		func() Iter[T] {
			elems := i.Next(ALL)
			met := make(map[T]struct{}, len(elems)*uniqRate)
			return Filter(func(t T) bool {
				if _, ok := met[t]; ok {
					return false
				}
				met[t] = struct{}{}
				return true
			}, StealSlice(elems))
		},
		func() Iter[T] {
			met := make(map[T]struct{})
			return Filter(func(t T) bool {
				if _, ok := met[t]; ok {
					return false
				}
				met[t] = struct{}{}
				return true
			}, i)
		},
	)
}

func UniqBy[T any, K comparable](f func(T) K, i Iter[T]) Iter[T] {
	return newFastpathIter(
		func() Iter[T] {
			elems := i.Next(ALL)
			met := make(map[K]struct{}, len(elems)*uniqRate)
			return Filter(func(t T) bool {
				k := f(t)
				if _, ok := met[k]; ok {
					return false
				}
				met[k] = struct{}{}
				return true
			}, StealSlice(elems))
		},
		func() Iter[T] {
			met := make(map[K]struct{})
			return Filter(func(t T) bool {
				k := f(t)
				if _, ok := met[k]; ok {
					return false
				}
				met[k] = struct{}{}
				return true
			}, i)
		},
	)
}

func Dup[T comparable](i Iter[T]) Iter[T] {
	return newFastpathIter(
		func() Iter[T] {
			elems := i.Next(ALL)
			met := make(map[T]bool, len(elems)*100/uniqRate)
			return Filter(func(t T) bool {
				return dupFilter(met, t)
			}, StealSlice(elems))
		},
		func() Iter[T] {
			met := make(map[T]bool)
			return Filter(func(t T) bool {
				return dupFilter(met, t)
			}, i)
		},
	)
}

func DupBy[T any, K comparable](f func(T) K, i Iter[T]) Iter[T] {
	return newFastpathIter(
		func() Iter[T] {
			elems := i.Next(ALL)
			met := make(map[K]bool, len(elems)*100/uniqRate)
			return Filter(func(t T) bool {
				return dupFilter(met, f(t))
			}, StealSlice(elems))
		},
		func() Iter[T] {
			met := make(map[K]bool)
			return Filter(func(t T) bool {
				return dupFilter(met, f(t))
			}, i)
		},
	)
}

func dupFilter[T comparable](met map[T]bool, k T) bool {
	isDup, ok := met[k]
	if ok && !isDup {
		met[k] = true
		return true
	}
	if !ok {
		met[k] = false
	}
	return false
}

func Sum[T gvalue.Numeric](i Iter[T]) T {
	var res T
	for _, e := range i.Next(ALL) {
		res += e
	}
	return res
}

func SumBy[T gvalue.Numeric](f func(T) T, i Iter[T]) T {
	var res T
	for _, e := range i.Next(ALL) {
		res += f(e)
	}
	return res
}

func Avg[T gvalue.Numeric](i Iter[T]) float64 {
	next := i.Next(ALL)
	if len(next) == 0 {
		return 0
	}

	var sum float64

	for _, e := range next {
		sum += float64(e)
	}
	return sum / float64(len(next))

}

func AvgBy[T gvalue.Numeric](f func(T) T, i Iter[T]) float64 {
	next := i.Next(ALL)
	if len(next) == 0 {
		return 0
	}

	var sum float64

	for _, e := range next {
		sum += float64(f(e))
	}
	return sum / float64(len(next))

}

func extend[T any](s *[]T, e []T) {
	if len(e) == 0 {
		return
	}
	if len(*s) == 0 {
		*s = e
		return
	}
	*s = append(*s, e...)
}

func GroupBy[T any, K comparable](f func(T) K, i Iter[T]) map[K][]T {
	m := make(map[K][]T)
	for _, e := range i.Next(ALL) {
		k := f(e)
		m[k] = append(m[k], e)
	}
	return m
}

func Remove[T comparable](v T, i Iter[T]) Iter[T] {
	return Filter(func(t T) bool {
		return t != v
	}, i)
}

func RemoveN[T comparable](v T, n int, i Iter[T]) Iter[T] {
	return Filter(func(t T) bool {
		if n <= 0 {
			return true
		}
		if v != t {
			return true
		}
		n--
		return false
	}, i)
}

type chunkIter[T any] struct {
	i Iter[T]
	n int
}

func (c *chunkIter[T]) Next(n int) [][]T {
	var next [][]T
	for n != 0 {
		v := c.i.Next(c.n)
		if len(v) == 0 {
			break
		}

		next = append(next, v)

		if len(v) != c.n {
			break
		}
		n--
	}
	return next
}

func Chunk[T any](n int, i Iter[T]) Iter[[]T] {
	return &chunkIter[T]{
		i: i,
		n: n,
	}
}

func Divide[T any](n int, i Iter[T]) [][]T {
	elems := i.Next(ALL)
	k := len(elems) / n
	m := len(elems) % n

	next := make([][]T, n)

	for idx := 0; idx < n; idx++ {
		next[idx] = elems[idx*k+min(idx, m) : (idx+1)*k+min(idx+1, m)]
	}
	return next
}

type compactIter[T comparable] struct {
	i    Iter[T]
	zero T
}

func (c *compactIter[T]) Next(n int) []T {
	next := c.i.Next(n)
	empty := n == ALL || len(next) < n

	n = len(next)

	j := 0

	for ; j < len(next); j++ {
		if next[j] == c.zero {
			break
		}
	}

	ptr := j

	if j < len(next) {
		todo := next[j+1:]
		for j < len(next) {
			for k := 0; k < len(todo); k++ {
				if todo[k] != c.zero {
					next[ptr] = todo[k]
					ptr++
				}
			}

			if empty || ptr == n {
				break
			}

			todo = c.i.Next(n - ptr)
			empty = len(todo) < n-ptr
		}
	}

	return next[:ptr]
}

func Compact[T comparable](i Iter[T]) Iter[T] {
	return &compactIter[T]{
		i: i,
	}
}
