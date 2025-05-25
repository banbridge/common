package gmap

import (
	"maps"

	"github.com/banbridge/common/pkg/gresult"
	"github.com/banbridge/common/pkg/gslice"
	"github.com/banbridge/common/pkg/gvalue"
	"github.com/banbridge/common/pkg/optional"
)

// Map 遍历map，对每个元素执行函数，返回一个新的map
// 注意：如果函数返回的key重复，则会覆盖前面的值
func Map[K1, K2 comparable, V1, V2 any](m map[K1]V1, f func(K1, V1) (K2, V2)) map[K2]V2 {
	ret := make(map[K2]V2, len(m))
	for k, v := range m {
		k2, v2 := f(k, v)
		ret[k2] = v2
	}
	return ret
}

// TryMap 遍历map，对每个元素执行函数，返回一个新的map
// 注意：如果函数返回的key重复，则会覆盖前面的值
func TryMap[K1, K2 comparable, V1, V2 any](m map[K1]V1, f func(K1, V1) (K2, V2, error)) gresult.R[map[K2]V2] {

	ret := make(map[K2]V2, len(m))
	for k, v := range m {
		k2, v2, err := f(k, v)
		if err != nil {
			return gresult.Err[map[K2]V2](err)
		}
		ret[k2] = v2
	}
	return gresult.OK(ret)
}

// MapKeys 遍历map，对每个元素执行函数，返回一个新的map
func MapKeys[K1, K2 comparable, V any](m map[K1]V, f func(K1) K2) map[K2]V {
	ret := make(map[K2]V, len(m))
	for k, v := range m {
		k2 := f(k)
		ret[k2] = v
	}
	return ret
}

// TryMapKeys 遍历map，对每个元素执行函数，返回一个新的map
func TryMapKeys[K1, K2 comparable, V any](m map[K1]V, f func(K1) (K2, error)) gresult.R[map[K2]V] {
	ret := make(map[K2]V, len(m))
	for k, v := range m {
		k2, err := f(k)
		if err != nil {
			return gresult.Err[map[K2]V](err)
		}
		ret[k2] = v
	}
	return gresult.OK(ret)
}

// MapValues 遍历map，对每个元素执行函数，返回一个新的map
// 注意：如果函数返回的key重复，则会覆盖前面的值
func MapValues[K comparable, V1, V2 any](m map[K]V1, f func(V1) V2) map[K]V2 {

	ret := make(map[K]V2, len(m))
	for k, v := range m {
		ret[k] = f(v)
	}
	return ret
}

// Filter 过滤map，返回一个新的map
func Filter[K comparable, V any](m map[K]V, f func(K, V) bool) map[K]V {
	ret := make(map[K]V, len(m)/2)
	for k, v := range m {
		if f(k, v) {
			ret[k] = v
		}
	}
	return ret
}

// FilterKeys 过滤map，返回一个新的map
func FilterKeys[K comparable, V any](m map[K]V, f func(K) bool) map[K]V {
	ret := make(map[K]V, len(m)/2)
	for k, v := range m {
		if f(k) {
			ret[k] = v
		}
	}
	return ret
}

// FilterByKeys 通过Key过滤map，返回一个新的map
func FilterByKeys[K comparable, V any](m map[K]V, keys ...K) map[K]V {
	ret := make(map[K]V, min(len(keys), len(m)))
	for _, k := range keys {
		if v, ok := m[k]; ok {
			ret[k] = v
		}
	}
	return ret
}

func FilterValues[K comparable, V any](m map[K]V, f func(V) bool) map[K]V {
	ret := make(map[K]V, len(m)/2)
	for k, v := range m {
		if f(v) {
			ret[k] = v
		}
	}
	return ret
}

func FilterByValues[K, V comparable](m map[K]V, values ...V) map[K]V {
	ret := make(map[K]V, min(len(values), len(m)))
	for k, v := range m {
		if gslice.Contains(values, v) {
			ret[k] = v
		}
	}
	return ret
}

// Reject 过滤map，返回一个新的map
// 注意：如果函数返回true，则会过滤掉该元素
// 注意：如果函数返回false，则会保留该元素
func Reject[K comparable, V any](m map[K]V, f func(K, V) bool) map[K]V {
	ret := make(map[K]V, len(m)/2)
	for k, v := range m {
		if !f(k, v) {
			ret[k] = v
		}
	}
	return ret
}

// RejectKeys 过滤map，返回一个新的map
// 注意：如果函数返回true，则会过滤掉该元素
// 注意：如果函数返回false，则会保留该元素
func RejectKeys[K comparable, V any](m map[K]V, f func(K) bool) map[K]V {
	ret := make(map[K]V, len(m)/2)
	for k, v := range m {
		if !f(k) {
			ret[k] = v
		}
	}
	return ret
}

// RejectByKeys 通过Key过滤map，返回一个新的map
func RejectByKeys[K comparable, V any](m map[K]V, keys ...K) map[K]V {
	ret := Clone(m)
	for k, v := range m {
		if !gslice.Contains(keys, k) {
			ret[k] = v
		}
	}
	return ret
}

func RejectValues[K comparable, V any](m map[K]V, f func(V) bool) map[K]V {
	ret := make(map[K]V, len(m)/2)
	for k, v := range m {
		if !f(v) {
			ret[k] = v
		}
	}
	return ret
}

func RejectByValues[K, V comparable](m map[K]V, values ...V) map[K]V {
	ret := make(map[K]V, len(m)/2)
	for k, v := range m {
		if !gslice.Contains(values, v) {
			ret[k] = v
		}
	}
	return ret
}

func Keys[K comparable, V any](m map[K]V) []K {
	if m == nil {
		return []K{}
	}
	ret := make([]K, 0, len(m))
	for k := range m {
		ret = append(ret, k)
	}
	return ret
}

func Values[K comparable, V any](m map[K]V) []V {
	if m == nil {
		return []V{}
	}
	ret := make([]V, 0, len(m))
	for _, v := range m {
		ret = append(ret, v)
	}
	return ret
}

func Merge[K comparable, V any](ms ...map[K]V) map[K]V {
	return Union(ms...)
}

func Union[K comparable, V any](ms ...map[K]V) map[K]V {
	if len(ms) == 0 {
		return make(map[K]V)
	}

	if len(ms) == 1 {
		return cloneWithoutCheck(ms[0])
	}

	var maxLen int
	for _, m := range ms {
		if len(m) > maxLen {
			maxLen = len(m)
		}
	}

	if maxLen == 0 {
		return make(map[K]V)
	}

	ret := make(map[K]V, maxLen)
	for _, m := range ms {
		for k, v := range m {
			ret[k] = v
		}
	}
	return ret
}

func UnionBy[K comparable, V any, M ~map[K]V](f ConflictFunc[K, V], ms ...map[K]V) map[K]V {
	if len(ms) == 0 {
		return make(map[K]V)
	}

	if len(ms) == 1 {
		return cloneWithoutCheck(ms[0])
	}

	var maxLen int
	for _, m := range ms {
		if len(m) > maxLen {
			maxLen = len(m)
		}
	}

	ret := make(map[K]V, maxLen)

	if maxLen == 0 {
		return ret
	}

	for _, m := range ms {
		for k, newVal := range m {
			if oldVal, ok := ret[k]; ok {
				ret[k] = f(k, oldVal, newVal)
			} else {
				ret[k] = newVal
			}
		}
	}
	return ret
}

func Diff[K comparable, V any](m map[K]V, against ...map[K]V) map[K]V {
	if len(m) == 0 {
		return make(map[K]V)
	}

	if len(against) == 0 {
		return cloneWithoutCheck(m)
	}

	ret := make(map[K]V, len(m))

	for k, v := range m {

		var ok bool
		for _, againstM := range against {
			if _, ok = againstM[k]; ok {
				break
			}
		}

		if !ok {
			ret[k] = v
		}

	}

	return ret
}

func Intersect[K comparable, V any](ms ...map[K]V) map[K]V {
	if len(ms) == 0 {
		return make(map[K]V)
	}

	if len(ms) == 1 {
		return cloneWithoutCheck(ms[0])
	}

	var minLen = len(ms[0])
	for _, m := range ms {
		minLen = min(minLen, len(m))
	}

	ret := make(map[K]V, minLen)
	if minLen == 0 {
		return ret
	}

	countMap := make(map[K]int, len(ms[0]))
	unionMap := make(map[K]V, len(ms[0]))
	for _, m := range ms {
		for k, v := range m {
			countMap[k]++
			unionMap[k] = v
		}
	}

	for k, v := range countMap {
		if v == len(ms) {
			ret[k] = unionMap[k]
		}
	}

	return ret
}

func IntersectBy[K comparable, V any, M ~map[K]V](f ConflictFunc[K, V], ms ...map[K]V) map[K]V {
	if len(ms) == 0 {
		return make(map[K]V)
	}
	if len(ms) == 1 {
		return cloneWithoutCheck(ms[0])
	}

	var minLen = len(ms[0])
	for _, m := range ms {
		minLen = min(minLen, len(m))
	}

	ret := make(map[K]V, minLen)
	if minLen == 0 {
		return ret
	}

	countMap := make(map[K]int, len(ms[0]))
	unionMap := make(map[K]V, len(ms[0]))
	for _, m := range ms {
		for k, v := range m {
			countMap[k]++
			if _, ok := unionMap[k]; ok {
				unionMap[k] = f(k, unionMap[k], v)
			}
		}
	}
	for k, v := range countMap {
		if v == len(ms) {
			ret[k] = f(k, unionMap[k], unionMap[k])
		}
	}
	return ret
}

func Load[K comparable, V any](m map[K]V, key K) optional.O[V] {
	if isEmpty(m) {
		return optional.Nil[V]()
	}

	v, ok := m[key]
	if !ok {
		return optional.Nil[V]()
	}
	return optional.OK(v)
}

func LoadOrStore[K comparable, V any](m map[K]V, key K, defaultV V) (v V, loaded bool) {
	if m == nil {
		m = make(map[K]V)
	}

	v, loaded = m[key]
	if !loaded {
		m[key] = defaultV
		v = defaultV
	}

	return
}

func LoadOrStoreFunc[K comparable, V any](m map[K]V, key K, defaultFunc func() V) (v V, loaded bool) {
	if m == nil {
		m = make(map[K]V)
	}
	v, loaded = m[key]
	if !loaded {
		v = defaultFunc()
		m[key] = v
	}
	return
}

func LoadAndDelete[K comparable, V any](m map[K]V, key K) optional.O[V] {
	if isEmpty(m) {
		return optional.Nil[V]()
	}
	v, ok := m[key]
	if !ok {
		return optional.Nil[V]()
	}
	delete(m, key)
	return optional.OK(v)
}

func LoadKey[K comparable, V comparable](m map[K]V, v V) optional.O[K] {
	if isEmpty(m) {
		return optional.Nil[K]()
	}
	for k, vv := range m {
		if vv == v {
			return optional.OK(k)
		}
	}
	return optional.Nil[K]()
}

func LoadBy[k comparable, V any](m map[k]V, f func(k k, v V) bool) optional.O[V] {
	if isEmpty(m) {
		return optional.Nil[V]()
	}
	for k, v := range m {
		if f(k, v) {
			return optional.OK(v)
		}
	}
	return optional.Nil[V]()
}

func LoadKeysBy[K comparable, V any](m map[K]V, f func(k K, v V) bool) optional.O[K] {
	if isEmpty(m) {
		return optional.Nil[K]()
	}

	for k, v := range m {
		if f(k, v) {
			return optional.OK(k)
		}
	}

	return optional.Nil[K]()
}

func Invert[K comparable, V comparable](m map[K]V) map[V]K {
	if isEmpty(m) {
		return make(map[V]K)
	}
	ret := make(map[V]K, len(m))
	for k, v := range m {
		ret[v] = k
	}
	return ret
}

func InvertBy[K comparable, V comparable, M ~map[K]V](m M, onConflict ConflictFunc[V, K]) map[V]K {
	if isEmpty(m) {
		return make(map[V]K)
	}

	ret := make(map[V]K, len(m))
	for k, v := range m {
		oldKey, ok := ret[v]
		if ok {
			ret[v] = onConflict(v, oldKey, k)
		} else {
			ret[v] = k
		}
	}
	return ret
}

func InvertGroup[K comparable, V comparable](m map[K]V) map[V][]K {
	if isEmpty(m) {
		return make(map[V][]K)
	}
	ret := make(map[V][]K, len(m))
	for k, v := range m {
		ret[v] = append(ret[v], k)
	}
	return ret
}

func Equal[K comparable, V comparable](m1 map[K]V, m2 map[K]V) bool {
	return maps.Equal(m1, m2)
}

func EqualBy[K comparable, V comparable, M ~map[K]V](m1 M, m2 M, eq func(V, V) bool) bool {
	return maps.EqualFunc(m1, m2, eq)
}

func EqualStrict[K comparable, V comparable](m1 map[K]V, m2 map[K]V) bool {
	if (m1 == nil && m2 != nil) || (m1 != nil && m2 == nil) {
		return false
	}
	return Equal(m1, m2)
}

func EqualStrictBy[K comparable, V comparable, M ~map[K]V](m1 M, m2 M, eq func(V, V) bool) bool {
	if (m1 == nil && m2 != nil) || (m1 != nil && m2 == nil) {
		return false
	}
	return EqualBy(m1, m2, eq)
}

type ConflictFunc[K comparable, V any] func(key K, oldVal V, newVal V) V

func DiscardOld[K comparable, V any]() ConflictFunc[K, V] {
	return discardOld[K, V]
}

func DiscardNew[K comparable, V any]() ConflictFunc[K, V] {
	return func(_ K, oldVal V, _ V) V {
		return oldVal
	}
}

func DiscardZero[K comparable, V comparable](fallback ConflictFunc[K, V]) ConflictFunc[K, V] {
	zeroVal := gvalue.Zero[V]()

	return func(k K, oldVal V, newVal V) V {
		if oldVal == zeroVal {
			if newVal != zeroVal {
				return newVal
			} else if fallback != nil {
				return fallback(k, oldVal, newVal)
			} else {
				return discardOld(k, oldVal, newVal)
			}
		} else {
			if newVal != zeroVal {
				return newVal
			} else if fallback != nil {
				return fallback(k, oldVal, newVal)
			} else {
				return discardOld(k, oldVal, newVal)
			}
		}
	}
}

func DiscardNil[K comparable, V any](fallback ConflictFunc[K, *V]) ConflictFunc[K, *V] {
	return func(k K, oldVal, newVal *V) *V {
		if oldVal == nil {
			if newVal != nil {
				return newVal
			} else if fallback != nil {
				return fallback(k, oldVal, newVal)
			} else {
				return discardOld(k, oldVal, newVal)
			}

		} else {
			if newVal != nil {
				return newVal
			} else if fallback != nil {
				return fallback(k, oldVal, newVal)
			} else {
				return discardOld(k, oldVal, newVal)
			}
		}

	}
}

func discardOld[K comparable, V any](_ K, _ V, newVal V) V {
	return newVal
}

func Clone[K comparable, V any](m map[K]V) map[K]V {
	if m == nil {
		return nil
	}
	return cloneWithoutCheck(m)
}

func CloneByKeys[K comparable, V any, M ~map[K]V](m M, f func(V) V) map[K]V {
	if m == nil {
		return nil
	}
	return MapValues(m, f)
}

func cloneWithoutCheck[K comparable, V any](m map[K]V) map[K]V {
	ret := make(map[K]V, len(m))
	for k, v := range m {
		ret[k] = v
	}
	return ret
}

func Contains[K comparable, V comparable](m map[K]V, key K) bool {
	if isEmpty(m) {
		return false
	}
	_, ok := m[key]

	return ok
}

func ContainsAny[K comparable, V comparable](m map[K]V, ks ...K) bool {
	if isEmpty(m) {
		return false
	}
	for _, k := range ks {
		if _, ok := m[k]; ok {
			return true
		}
	}
	return false
}

func ContainsAll[K comparable, V comparable](m map[K]V, ks ...K) bool {
	if isEmpty(m) {
		return false
	}
	for _, k := range ks {
		if _, ok := m[k]; !ok {
			return false
		}
	}
	return true
}

func LoadAll[K comparable, V any](m map[K]V, ks ...K) []V {
	if isEmpty(m) {
		return nil
	}
	ret := make([]V, 0, len(ks))
	for _, k := range ks {
		if v, ok := m[k]; ok {
			ret = append(ret, v)
		} else {
			ret = append(ret, v)
		}
	}
	return ret
}

func LoadAny[K comparable, V any](m map[K]V, ks ...K) (r optional.O[V]) {
	if isEmpty(m) || len(ks) == 0 {
		return
	}
	for _, k := range ks {
		if v, ok := m[k]; ok {
			return optional.OK(v)
		}
	}

	return

}

func LoadSome[K comparable, V any](m map[K]V, ks ...K) (r []V) {
	if isEmpty(m) || len(ks) == 0 {
		return
	}
	for _, k := range ks {
		if v, ok := m[k]; ok {
			r = append(r, v)
		}
	}
	return
}

func Sum[K comparable, V gvalue.Numeric](m map[K]V) V {
	var ret V
	for _, v := range m {
		ret += v
	}

	return ret
}

func SumBy[K comparable, V any, N gvalue.Numeric](m map[K]V, f func(V) N) N {
	var ret N
	for _, v := range m {
		ret += f(v)
	}
	return ret
}

func Avg[K comparable, V gvalue.Numeric](m map[K]V) float64 {
	var ret V
	var count int
	for _, v := range m {
		ret += v
		count++
	}
	if count == 0 {
		return 0
	}
	return float64(ret) / float64(count)
}

func AvgBy[K comparable, V any, N gvalue.Numeric](m map[K]V, f func(V) N) float64 {
	var ret N
	var count int
	for _, v := range m {
		ret += f(v)
		count++
	}
	if count == 0 {
		return 0
	}
	return float64(ret) / float64(count)
}

func Max[K comparable, V gvalue.Ordered](m map[K]V) (r optional.O[V]) {

	if isEmpty(m) {
		return optional.Nil[V]()
	}

	var maxV V

	for _, v := range m {
		if v > maxV {
			maxV = v
		}
	}
	return optional.OK(maxV)
}

func isEmpty[K comparable, V any](m map[K]V) bool {
	return m == nil || len(m) == 0
}
