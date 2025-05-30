package gslice

import (
	"sort"

	"slices"

	"github.com/banbridge/common/pkg/collection/set"
	"github.com/banbridge/common/pkg/gptr"
	"github.com/banbridge/common/pkg/gvalue"
)

func ForEach[T any](list []T, eachFunc func(item T)) {
	for _, val := range list {
		eachFunc(val)
	}
}

// Map 循环遍历切片，对每个元素执行函数，返回一个新的切片
func Map[T any, R any](list []T, mapFunc func(item T) R) []R {
	result := make([]R, len(list))
	for i, val := range list {
		result[i] = mapFunc(val)
	}
	return result
}

// TryMap 循环遍历切片，对每个元素执行函数，如果函数返回错误，则返回错误
func TryMap[T, R any](s []T, f func(T) (R, error)) ([]R, error) {
	var result []R
	for _, v := range s {
		r, err := f(v)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}

// Filter 循环遍历切片，对每个元素执行函数，如果函数返回true，则添加到新的切片中
func Filter[T any](list []T, filterFunc func(item T) bool) []T {
	result := make([]T, 0)
	for _, val := range list {
		if filterFunc(val) {
			result = append(result, val)
		}
	}
	return result
}

// FilterMap 循环遍历切片，对每个元素执行函数，如果函数返回true，则添加到新的切片中
func FilterMap[T any, R any](list []T, filterFunc func(item T) (R, bool)) []R {
	result := make([]R, 0)
	for _, val := range list {
		r, ok := filterFunc(val)
		if ok {
			result = append(result, r)
		}
	}
	return result
}

// TryFilterMap 循环遍历切片，对每个元素执行函数，如果函数返回true，则添加到新的切片中，如果函数返回false，则添加到notMatch切片中
func TryFilterMap[T any, R any](list []T, filterFunc func(item T) (R, bool)) ([]R, []T) {
	result := make([]R, 0)
	notMatch := make([]T, 0)
	for _, val := range list {
		r, ok := filterFunc(val)
		if ok {
			result = append(result, r)
		} else {
			notMatch = append(notMatch, val)
		}
	}
	return result, notMatch
}

func MapInplace[T any](list []T, mapFunc func(item T) T) []T {
	for i, val := range list {
		list[i] = mapFunc(val)
	}
	return list
}

// Reject 循环遍历切片，对每个元素执行函数，如果函数返回false，则添加到新的切片中
func Reject[T any](list []T, filterFunc func(item T) bool) []T {
	result := make([]T, 0, len(list))
	for _, val := range list {
		if !filterFunc(val) {
			result = append(result, val)
		}
	}
	return result
}

// AnyMatch 循环遍历切片，对每个元素执行函数，如果函数返回true，则返回true
func AnyMatch[T any](list []T, matchFunc func(item T) bool) bool {
	for _, val := range list {
		if matchFunc(val) {
			return true
		}
	}
	return false
}

// AllMatch 循环遍历切片，对每个元素执行函数，如果函数返回false，则返回false
func AllMatch[T any](list []T, matchFunc func(item T) bool) bool {
	for _, val := range list {
		if !matchFunc(val) {
			return false
		}
	}
	return true
}

// MapReduce map后reduce
func MapReduce[T, M, R any](list []T, mapFunc func(T) M, reduceFunc func() func(M) R) R {
	var res R
	for _, val := range list {
		res = reduceFunc()(mapFunc(val))
	}
	return res
}

// GroupBy 分组
func GroupBy[T any, K comparable](list []T, groupFunc func(T) K) map[K][]T {
	res := make(map[K][]T)
	for _, val := range list {
		key := groupFunc(val)

		mList, ok := res[key]
		if !ok {
			res[key] = make([]T, 0)
		}
		mList = append(mList, val)
		res[key] = mList
	}
	return res
}

// Contains 切片是否包含元素
func Contains[T comparable](list []T, item T) bool {
	for _, val := range list {
		if val == item {
			return true
		}
	}
	return false
}

// ContainsAny 切片是否包含任意元素
func ContainsAny[T comparable](list []T, items ...T) bool {
	m := make(map[T]struct{}, len(items))
	for _, item := range items {
		m[item] = struct{}{}
	}
	for _, val := range list {
		if _, ok := m[val]; ok {
			return true
		}
	}
	return false
}

// ContainsAll 切片是否包含所有元素
func ContainsAll[T comparable](list []T, items ...T) bool {
	m := make(map[T]struct{}, len(items))
	for _, item := range items {
		m[item] = struct{}{}
	}

	for _, val := range list {
		delete(m, val)
		if len(m) == 0 {
			return true
		}
	}
	return len(m) == 0
}

// Uniq 去重
func Uniq[T comparable](list []T) []T {
	res := make([]T, 0, len(list))
	m := make(map[T]struct{}, len(list))
	for _, val := range list {
		if _, ok := m[val]; ok {
			continue
		}
		m[val] = struct{}{}
		res = append(res, val)
	}
	return res
}

// UniqBy 去重
func UniqBy[T any, K comparable](list []T, f func(T) K) []T {
	res := make([]T, 0, len(list))
	m := make(map[K]struct{}, len(list))
	for _, val := range list {
		key := f(val)
		if _, ok := m[key]; ok {
			continue
		}
		m[key] = struct{}{}
		res = append(res, val)
	}
	return res
}

// Duplicate 找出重复的元素
func Duplicate[T comparable](list []T) []T {
	res := make([]T, 0, len(list))
	m := make(map[T]struct{}, len(list))
	for _, val := range list {
		if _, ok := m[val]; ok {
			res = append(res, val)
			continue
		}
		m[val] = struct{}{}
	}
	return res
}

// DuplicateBy 找出重复的元素
func DuplicateBy[T any, K comparable](list []T, f func(T) K) []T {
	res := make([]T, 0, len(list))
	m := make(map[K]struct{}, len(list))
	for _, val := range list {
		key := f(val)
		if _, ok := m[key]; ok {
			res = append(res, val)
			continue
		}
		m[key] = struct{}{}
	}
	return res
}

// Repeat 重复元素
func Repeat[T any](item T, n int) []T {
	if n < 0 {
		panic("repeat count is negative")
	}

	res := make([]T, n)
	for i := 0; i < n; i++ {
		res[i] = item
	}
	return res
}

// RepeatBy 重复元素
func RepeatBy[T any](f func() T, n int) []T {
	if n < 0 {
		panic("repeat count is negative")
	}

	res := make([]T, n)
	for i := 0; i < n; i++ {
		res[i] = f()
	}
	return res
}

func Clone[T any](list []T) []T {
	return slices.Clone(list)
}

func Flatten[T any](list [][]T) []T {
	res := make([]T, 0, len(list))
	for _, val := range list {
		res = append(res, val...)
	}
	return res
}

func Union[T comparable](ss ...[]T) []T {
	if len(ss) == 0 {
		return []T{}
	}
	if len(ss) == 1 {
		return Uniq(ss[0])
	}

	members := set.New[T]()
	res := make([]T, 0, 32)

	for _, s := range ss {
		for _, member := range s {
			if members.Add(member) {
				res = append(res, member)
			}
		}
	}
	return res
}

func Diff[T comparable](s []T, against ...[]T) []T {
	if len(against) == 0 {
		return s
	}
	if len(against) == 0 {
		return Uniq(s)
	}

	members := set.New[T](s...)

	for _, tmpList1 := range against {
		for _, member := range tmpList1 {
			members.Remove(member)
		}
	}

	if members.Len() == 0 {
		return []T{}
	}

	res := make([]T, 0, members.Len())

	for _, v := range s {
		if members.Remove(v) {
			res = append(res, v)
			if members.Len() == 0 {
				return res
			}
		}
	}

	return res
}

func Intersect[T comparable](ss ...[]T) []T {
	if len(ss) == 0 {
		return []T{}
	}

	if len(ss) == 1 {
		return Uniq(ss[0])
	}

	if len(ss[0]) == 0 {
		return []T{}
	}

	members := set.New[T](ss[0]...)

	for _, tmpList := range ss[1:] {
		if len(tmpList) == 0 {
			return []T{}
		}
		members.IntersectInPlace(set.New[T](tmpList...))
	}

	if members.Len() == 0 {
		return []T{}
	}

	res := make([]T, 0, members.Len())

	for _, tmpList := range ss {
		for _, member := range tmpList {
			if members.Remove(member) {
				res = append(res, member)
				if members.Len() == 0 {
					return res
				}
			}
		}
	}
	return res
}

// Sum 求和
func Sum[T gvalue.Numeric](list []T) T {
	var res T
	for _, val := range list {
		res += val
	}
	return res
}

func SumBy[T any, K gvalue.Numeric](list []T, f func(T) K) K {
	var res K
	for _, val := range list {
		res += f(val)
	}
	return res
}

// Avg 求平均值
func Avg[T gvalue.Numeric](list []T) float64 {
	if len(list) == 0 {
		return 0
	}
	var res float64
	for _, val := range list {
		res += float64(val)
	}

	return res / float64(len(list))
}

// AvgBy 求平均值
func AvgBy[T any, K gvalue.Numeric](list []T, f func(T) K) float64 {
	if len(list) == 0 {
		return 0
	}
	var res float64
	for _, val := range list {
		res += float64(f(val))
	}
	return res / float64(len(list))
}

// Compact 去除零值
func Compact[T comparable](list []T) []T {
	return Filter(list, gvalue.IsNotZero)
}

func normalizeIndex[T any, I gvalue.Integer](list []T, index I) (int, bool) {
	m := int(index)
	if m < 0 {
		m += len(list)
	}
	return m, m > 0 && m < len(list)
}

func Slice[T any, I gvalue.Integer](list []T, start, end I) []T {
	startIdx, _ := normalizeIndex(list, start)

	var endIdx int
	if start < 0 && end == 0 {
		endIdx = len(list)
	} else {
		endIdx, _ = normalizeIndex(list, end)
	}

	if startIdx < 0 {
		startIdx = 0
	}

	if endIdx > len(list) {
		endIdx = len(list)
	}

	if startIdx > endIdx {
		return []T{}
	}

	return list[start:endIdx]
}

func SliceClone[T any, I gvalue.Integer](list []T, start, end I) []T {
	return Clone(Slice(list, start, end))
}

func Of[T any](v ...T) []T {
	if len(v) == 0 {
		return []T{}
	}
	return v
}

func Count[T comparable](list []T, item T) int {
	var count int
	for _, val := range list {
		if val == item {
			count++
		}
	}
	return count
}

func CountBy[T any](list []T, f func(T) bool) int {
	var count int
	for _, val := range list {
		if f(val) {
			count++
		}
	}
	return count
}

// CountValues 统计每个元素出现的次数
func CountValues[T comparable](list []T) map[T]int {
	res := make(map[T]int, len(list)/2)
	for i := range list {
		res[list[i]]++
	}
	return res
}

// CountValuesBy 统计每个元素出现的次数
func CountValuesBy[T any, K comparable](list []T, f func(T) K) map[K]int {
	res := make(map[K]int, len(list)/2)
	for i := range list {
		res[f(list[i])]++
	}
	return res
}

// Partition 分割切片，将切片中的元素按照函数的返回值分成两个切片
func Partition[T any](list []T, f func(T) bool) ([]T, []T) {
	trueList := make([]T, 0, len(list)/2)
	falseList := make([]T, 0, len(list)/2)
	for i := range list {
		if f(list[i]) {
			trueList = append(trueList, list[i])
		} else {
			falseList = append(falseList, list[i])
		}
	}
	return trueList, falseList
}

// Remove 移除元素
func Remove[T comparable](list []T, item T) []T {
	for i, val := range list {
		if val == item {
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}

// Fold 折叠切片，将切片中的元素依次传入函数，返回最终的结果
func Fold[T, R any](list []T, f func(R, T) R, init R) R {
	for _, val := range list {
		init = f(init, val)
	}
	return init
}

// Reverse 反转切片
func Reverse[T any](list []T) []T {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
	return list
}

// All 循环遍历切片，对每个元素执行函数，如果函数返回false，则返回false
// 注意：如果切片为空，返回true
func All[T any](list []T, f func(T) bool) bool {
	for _, val := range list {
		if !f(val) {
			return false
		}
	}
	return true
}

// Any 循环遍历切片，对每个元素执行函数
// 注意：如果切片为空，返回false
func Any[T any](list []T, f func(T) bool) bool {
	for _, val := range list {
		if f(val) {
			return true
		}
	}
	return false
}

// And 循环遍历切片，对每个元素执行函数，
func And[T ~bool](list []T) bool {
	for _, val := range list {
		if !val {
			return false
		}
	}
	return true
}

func Or[T ~bool](list []T) bool {
	for _, val := range list {
		if val {
			return true
		}
	}
	return false
}

type sortable[T any] struct {
	items    []T
	lessFunc func(T, T) bool
}

func (s *sortable[T]) Len() int {
	return len(s.items)
}

func (s *sortable[T]) Less(i, j int) bool {
	return s.lessFunc(s.items[i], s.items[j])
}

func (s *sortable[T]) Swap(i, j int) {
	tmp := s.items[i]
	s.items[i] = s.items[j]
	s.items[j] = tmp
}

func Sort[T any](items []T, lessFunc func(T, T) bool) []T {
	s := &sortable[T]{
		items:    items,
		lessFunc: lessFunc,
	}
	sort.Sort(s)
	return items
}

func StableSort[T any](items []T, lessFunc func(T, T) bool) []T {
	s := &sortable[T]{
		items:    items,
		lessFunc: lessFunc,
	}
	sort.Stable(s)
	return items
}

func ForEachIndex[T any](items []T, f func(int, T)) {
	for i, val := range items {
		f(i, val)
	}
}

func Equal[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func EqualBy[T any, K comparable](a, b []T, f func(T) K) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if f(v) != f(b[i]) {
			return false
		}
	}
	return true
}

func ToMapValues[T any, K comparable](list []T, f func(T) K) map[K]T {
	res := make(map[K]T, len(list))
	for _, val := range list {
		res[f(val)] = val
	}
	return res
}

func ToMap[T, V any, K comparable](list []T, f func(T) (K, V)) map[K]V {

	res := make(map[K]V, len(list))
	for _, val := range list {
		k, v := f(val)
		res[k] = v
	}
	return res
}

// Divide 将切片分成n份
// 如果n大于切片长度，则返回一个包含切片的切片
// 如果n小于等于0，则返回一个空切片
func Divide[T any](list []T, n int) [][]T {
	k := len(list) / n
	m := len(list) % n

	res := make([][]T, n)

	for i := 0; i < n; i++ {
		res[i] = list[i*k+min(i, m) : (i+1)*k+min(i+1, m)]
	}

	return res
}

func PtrOf[T any](list []T) []*T {
	return Map(list, gptr.Of[T])
}

func IndirectOf[T any](list []*T) []T {
	return Map(list, gptr.Indirect[T])
}

func Index[T comparable](list []T, item T) int {
	return slices.Index(list, item)
}

func IndexRev[T comparable](list []T, item T) int {
	for i := len(list) - 1; i >= 0; i-- {
		if list[i] == item {
			return i
		}
	}
	return -1
}

func IndexBy[T any, K comparable](list []T, f func(T) bool) int {
	return slices.IndexFunc(list, f)
}

func IndexByRev[T any, K comparable](list []T, f func(T) bool) int {
	for i := len(list) - 1; i >= 0; i-- {
		if f(list[i]) {
			return i
		}
	}
	return -1
}
