package metrics

import (
	"sync/atomic"

	"golang.org/x/exp/constraints"
)

type GenericGauge[T constraints.Ordered] interface {
	Snapshot() GenericGauge[T]
	Update(T)
	Value() T
}

func GetOrRegisterGenericGauge[T constraints.Ordered](name string, r Registry) GenericGauge[T] {
	if r == nil {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewGenericGauge[T]).(GenericGauge[T])
}

func NewGenericGauge[T constraints.Ordered]() GenericGauge[T] {
	if UseNilMetrics {
		return NilGenericGauge[T]{}
	}
	v := atomic.Pointer[T]{}
	var zero T
	v.Store(&zero)
	return &StandardGenericGauge[T]{value: v}
}

func NewRegisteredGenericGauge[T constraints.Ordered](name string, r Registry) GenericGauge[T] {
	g := NewGenericGauge[T]()
	if r == nil {
		r = DefaultRegistry
	}
	r.Register(name, g)
	return g
}

// mimicking `type GaugeSnapshot int64`
type GenericGaugeSnapshot[T constraints.Ordered] struct{ X T }

func (GenericGaugeSnapshot[T]) Update(T) {
	panic("Update called on a GenericGaugeSnapshot")
}

func (g GenericGaugeSnapshot[T]) Value() T {
	return g.X
}

func (g GenericGaugeSnapshot[T]) Snapshot() GenericGauge[T] {
	return g
}

type NilGenericGauge[T constraints.Ordered] struct{}

func (NilGenericGauge[T]) Snapshot() GenericGauge[T] {
	return NilGenericGauge[T]{}
}

func (NilGenericGauge[T]) Update(v T) {

}

func (NilGenericGauge[T]) Value() T {
	var zero T
	return zero
}

type StandardGenericGauge[T constraints.Ordered] struct {
	value atomic.Pointer[T]
}

func (g *StandardGenericGauge[T]) Value() T {
	// Do we actually need atomic load if we're being fronted by channels?
	return *g.value.Load()
}
func (g *StandardGenericGauge[T]) Update(v T) {
	// Do we actually need atomic store if we're being fronted by channels?
	g.value.Store(&v)
}
func (g *StandardGenericGauge[T]) Snapshot() GenericGauge[T] {
	return g
}

type FunctionalGenericGauge[T constraints.Ordered] struct {
	value func() T
}

func (g FunctionalGenericGauge[T]) Value() T {
	return g.value()
}
func (g FunctionalGenericGauge[T]) Update(T) {
	panic("Update called on a FunctionalGenericGauge")
}
func (g FunctionalGenericGauge[T]) Snapshot() GenericGauge[T] {
	return GenericGauge[T](g)
}
