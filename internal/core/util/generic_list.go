package util

import "errors"

type GenericList[T comparable] struct {
	Items []T
}

func NewGenericList[T comparable]() *GenericList[T] {
	return &GenericList[T]{}

}

func (gl *GenericList[T]) Add(item T) {
	gl.Items = append(gl.Items, item)
}

func (gl *GenericList[T]) Remove(index int) error {
	if index < 0 || index >= len(gl.Items) {
		return errors.New("index out of range")
	}
	gl.Items = append(gl.Items[:index], gl.Items[index+1:]...)
	return nil
}

func (gl *GenericList[T]) RemoveWithValue(item T) {
	for i, v := range gl.Items {
		if v == item {
			gl.Items = append(gl.Items[:i], gl.Items[i+1:]...)
			break
		}
	}
}

func (gl *GenericList[T]) Contains(item T) bool {
	for _, v := range gl.Items {
		if v == item {
			return true
		}
	}
	return false
}

func (gl *GenericList[T]) Get(index int) (T, error) {
	if index < 0 || index >= len(gl.Items) {
		var zeroValue T
		return zeroValue, errors.New("index out of range")
	}
	return gl.Items[index], nil
}
