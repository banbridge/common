package stream

// Each 循环遍历切片，对每个元素执行函数，如果函数返回错误，则停止循环
func Each[T any](list []T, eachFunc func(item T) error) error {
	for _, val := range list {
		err := eachFunc(val)
		if err != nil {
			return err
		}
	}
	return nil
}

// Map 循环遍历切片，对每个元素执行函数，返回一个新的切片
func Map[T any, R any](list []T, mapFunc func(item T) R) []R {
	result := make([]R, len(list))
	for i, val := range list {
		result[i] = mapFunc(val)
	}
	return result
}

// MapWithError 循环遍历切片，对每个元素执行函数，返回一个新的切片，如果函数返回错误，则停止循环
func MapWithError[T any, R any](list []T, mapFunc func(item T) (R, error)) ([]R, error) {
	result := make([]R, len(list))
	for i, val := range list {
		r, err := mapFunc(val)
		if err != nil {
			return nil, err
		}
		result[i] = r
	}
	return result, nil
}

// MapWithFilter 循环遍历切片，对每个元素执行函数，返回一个新的切片，如果函数返回false，则不添加到新的切片中
func MapWithFilter[T any, R any](list []T, mapFunc func(item T) (R, bool)) []R {
	result := make([]R, 0)
	for _, val := range list {
		r, ok := mapFunc(val)
		if ok {
			result = append(result, r)
		}
	}
	return result
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

// FilterWithError 循环遍历切片，对每个元素执行函数，如果函数返回错误，则停止循环
func FilterWithError[T any](list []T, filterFunc func(item T) (bool, error)) ([]T, error) {
	result := make([]T, 0)
	for _, val := range list {
		r, err := filterFunc(val)
		if err != nil {
			return nil, err
		}
		if r {
			result = append(result, val)
		}
	}
	return result, nil
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

// SliceToMap 切片转map
func SliceToMap[T any, K comparable](list []T, keyFunc func(T) K) map[K]T {
	res := make(map[K]T)
	for _, val := range list {
		key := keyFunc(val)
		res[key] = val
	}
	return res
}

// MapValues map值转切片
func MapValues[K comparable, V any](m map[K]V) []V {
	res := make([]V, 0)
	for _, val := range m {
		res = append(res, val)
	}
	return res
}
