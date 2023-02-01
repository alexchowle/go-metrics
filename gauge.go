package metrics

import (
	"sync/atomic"
)

// Gauges hold an int64 value that can be set arbitrarily.
type Gauge interface {
	Snapshot() Gauge
	Update(int64)
	Value() int64
}

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

// GetOrRegisterGauge returns an existing Gauge or constructs and registers a
// new StandardGauge.
func GetOrRegisterGauge(name string, r Registry) Gauge {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewGauge).(Gauge)
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

// NewGauge constructs a new StandardGauge.
func NewGauge() Gauge {
	if UseNilMetrics {
		return NilGauge{}
	}
	return &StandardGauge{0}
}

func NewRegisteredGenericGauge[T Number](name string, r Registry) GenericGauge[T] {
	g := NewGenericGauge[T]()
	if r == nil {
		r = DefaultRegistry
	}
	r.Register(name, g)
	return g
}

// NewRegisteredGauge constructs and registers a new StandardGauge.
func NewRegisteredGauge(name string, r Registry) Gauge {
	c := NewGauge()
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// NewFunctionalGauge constructs a new FunctionalGauge.
func NewFunctionalGauge(f func() int64) Gauge {
	if UseNilMetrics {
		return NilGauge{}
	}
	return &FunctionalGauge{value: f}
}

// NewRegisteredFunctionalGauge constructs and registers a new StandardGauge.
func NewRegisteredFunctionalGauge(name string, r Registry, f func() int64) Gauge {
	c := NewFunctionalGauge(f)
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// mimicking `type GaugeSnapshot int64`
type GenericGaugeSnapshot[T Number] struct{ X T }

func (GenericGaugeSnapshot[T]) Update(T) {
	panic("Update called on a GenericGaugeSnapshot")
}
func (g GenericGaugeSnapshot[T]) Value() T                  { return g.X }
func (g GenericGaugeSnapshot[T]) Snapshot() GenericGauge[T] { return g }

// GaugeSnapshot is a read-only copy of another Gauge.
type GaugeSnapshot int64

// Snapshot returns the snapshot.
func (g GaugeSnapshot) Snapshot() Gauge { return g }

// Update panics.
func (GaugeSnapshot) Update(int64) {
	panic("Update called on a GaugeSnapshot")
}

// Value returns the value at the time the snapshot was taken.
func (g GaugeSnapshot) Value() int64 { return int64(g) }

type NilGenericGauge[T Number] struct{}

func (NilGenericGauge[T]) Snapshot() GenericGauge[T] { return NilGenericGauge[T]{} }
func (NilGenericGauge[T]) Update(v T)                {}
func (NilGenericGauge[T]) Value() T                  { return 0 }

// NilGauge is a no-op Gauge.
type NilGauge struct{}

// Snapshot is a no-op.
func (NilGauge) Snapshot() Gauge { return NilGauge{} }

// Update is a no-op.
func (NilGauge) Update(v int64) {}

// Value is a no-op.
func (NilGauge) Value() int64 { return 0 }

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

// StandardGauge is the standard implementation of a Gauge and uses the
// sync/atomic package to manage a single int64 value.
type StandardGauge struct {
	value int64
}

// Snapshot returns a read-only copy of the gauge.
func (g *StandardGauge) Snapshot() Gauge {
	return GaugeSnapshot(g.Value())
}

// Update updates the gauge's value.
func (g *StandardGauge) Update(v int64) {
	atomic.StoreInt64(&g.value, v)
}

// Value returns the gauge's current value.
func (g *StandardGauge) Value() int64 {
	return atomic.LoadInt64(&g.value)
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
func (g FunctionalGenericGauge[T]) Snapshot() GenericGauge[T] { return GenericGauge[T](g) }

// FunctionalGauge returns value from given function
type FunctionalGauge struct {
	value func() int64
}

// Value returns the gauge's current value.
func (g FunctionalGauge) Value() int64 {
	return g.value()
}

// Snapshot returns the snapshot.
func (g FunctionalGauge) Snapshot() Gauge { return GaugeSnapshot(g.Value()) }

// Update panics.
func (FunctionalGauge) Update(int64) {
	panic("Update called on a FunctionalGauge")
}
