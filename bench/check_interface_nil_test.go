package bench

import (
	"reflect"
	"testing"
)

type Interface interface{}
type concreteType struct{}

func getConcreteType() *concreteType {
	return nil
}

func checkNil(v interface{}) {
	if v != nil {
	}
}

func checkNilWithReflect(v interface{}) {
	value := reflect.ValueOf(v)
	if !value.IsNil() {
	}
}

func BenchmarkCheckNilWithoutReflectInInterface(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d := getConcreteType()
		checkNil(d)
	}
}

func BenchmarkCheckNilWithReflectInInterface(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d := getConcreteType()
		checkNilWithReflect(d)
	}
}
