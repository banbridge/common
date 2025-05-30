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

// MapValues map值转切片
func MapValues[K comparable, V any](m map[K]V) []V {
	res := make([]V, 0)
	for _, val := range m {
		res = append(res, val)
	}
	return res
}
