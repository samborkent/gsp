package gsp

// Multi-channel type.
//
//	len(*new(MultiChannel[T])) == numChannels
type MultiChannel[T Type] []T

func (c MultiChannel[T]) Add(x T) MultiChannel[T] {
	for i := range c {
		c[i] += x
	}

	return c
}

func (c MultiChannel[T]) AddMC(x MultiChannel[T]) MultiChannel[T] {
	if len(c) != len(x) {
		return c
	}

	for i := range c {
		c[i] += x[i]
	}

	return c
}

func (c MultiChannel[T]) Divide(x T) MultiChannel[T] {
	for i := range c {
		c[i] /= x
	}

	return c
}

func (c MultiChannel[T]) DivideMC(x MultiChannel[T]) MultiChannel[T] {
	if len(c) != len(x) {
		return c
	}

	for i := range c {
		c[i] /= x[i]
	}

	return c
}

// M returns the mid channel.
func (c MultiChannel[T]) M() T {
	sum := T(0)

	for _, sample := range c {
		sum += sample
	}

	sum /= T(len(c))

	return sum
}

func (c MultiChannel[T]) Multiply(x T) MultiChannel[T] {
	for i := range c {
		c[i] *= x
	}

	return c
}

func (c MultiChannel[T]) MultiplyMC(x MultiChannel[T]) MultiChannel[T] {
	if len(c) != len(x) {
		return c
	}

	for i := range c {
		c[i] *= x[i]
	}

	return c
}

func (c MultiChannel[T]) Set(x T) MultiChannel[T] {
	for i := range c {
		c[i] = x
	}

	return c
}

func (c MultiChannel[T]) Subtract(x T) MultiChannel[T] {
	for i := range c {
		c[i] -= x
	}

	return c
}

func (c MultiChannel[T]) SubtractMC(x MultiChannel[T]) MultiChannel[T] {
	if len(c) != len(x) {
		return c
	}

	for i := range c {
		c[i] -= x[i]
	}

	return c
}

func ToMultiChannel[T Type](samples ...T) MultiChannel[T] {
	c := make(MultiChannel[T], len(samples))
	copy(c, samples)
	return c
}

func ZeroMultiChannel[T Type](n int) MultiChannel[T] {
	c := make(MultiChannel[T], n)

	switch any(T(0)).(type) {
	case uint8, uint16, uint32, uint64:
		c.Set(Zero[T]())
	}

	return c
}
