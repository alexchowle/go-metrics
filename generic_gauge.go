package metrics

import "sync/atomic"

type Number interface {
	int64 | float64
}

type GenericGauge[T Number] interface {
	Snapshot() GenericGauge[T]
	Update(T)
	Value() T
}

func GetOrRegisterGenericGauge[T Number](name string, r Registry) GenericGauge[T] {
	if r == nil {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewGenericGauge[T]).(GenericGauge[T])
}

func NewGenericGauge[T Number]() GenericGauge[T] {
	if UseNilMetrics {
		return NilGenericGauge[T]{}
	}
	v := atomic.Pointer[T]{}
	var zero T = 0
	v.Store(&zero)
	return &StandardGenericGauge[T]{value: v}
}

func NewRegisteredGenericGauge[T Number](name string, r Registry) GenericGauge[T] {
	g := NewGenericGauge[T]()
	if r == nil {
		r = DefaultRegistry
	}
	r.Register(name, g)
	return g
}

// mimicking `type GaugeSnapshot int64`
type GenericGaugeSnapshot[T Number] struct{ X T }

func (GenericGaugeSnapshot[T]) Update(T) {
	panic("Update called on a GenericGaugeSnapshot")
}

func (g GenericGaugeSnapshot[T]) Value() T {
	return g.X
}

func (g GenericGaugeSnapshot[T]) Snapshot() GenericGauge[T] {
	return g
}

type NilGenericGauge[T Number] struct{}

func (NilGenericGauge[T]) Snapshot() GenericGauge[T] {
	return NilGenericGauge[T]{}
}

func (NilGenericGauge[T]) Update(v T) {

}

func (NilGenericGauge[T]) Value() T {
	return 0
}

type StandardGenericGauge[T Number] struct {
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

type FunctionalGenericGauge[T Number] struct {
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
