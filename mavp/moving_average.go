// Package mavp provides a generic implementation of a moving average with a
// variable period.
//
// The intended use case is real time data visualization, where it might be
// desirable to change the period if the frequency at which data is being
// received changes.
package mavp

import (
	"github.com/jmaralo/go-utils/generics"
	"github.com/jmaralo/go-utils/rinbuf"
)

type MovingAverage[N generics.Float] struct {
	buffer   rinbuf.RingBuffer[N]
	elements int
	currAvg  N
}

// Create a new MovingAverage with the given period.
func New[N generics.Float](period int) MovingAverage[N] {
	return MovingAverage[N]{
		buffer:   rinbuf.New[N](period),
		elements: 0,
		currAvg:  0,
	}
}

// Add a new value to the moving average. Returns the current average.
func (avg *MovingAverage[N]) Add(value N) (current N) {
	if avg.elements >= avg.buffer.Len() {
		return avg.addValue(value)
	}

	return avg.addElement(value)
}

func (avg *MovingAverage[N]) addValue(value N) (current N) {
	avg.currAvg -= avg.buffer.Push(value) / N(avg.buffer.Len())
	avg.currAvg += value / N(avg.buffer.Len())
	return avg.currAvg
}

func (avg *MovingAverage[N]) addElement(value N) (current N) {
	avg.buffer.Push(value)
	avg.elements += 1
	avg.currAvg += (value - avg.currAvg) / N(avg.elements)
	return avg.currAvg
}

// Returns the current average.
func (avg *MovingAverage[N]) Current() N {
	return avg.currAvg
}

// Change the period of the moving average. If the period is less
// or equal to zero this function returns an error
func (avg *MovingAverage[N]) Resize(period int) error {
	if period > avg.buffer.Len() {
		avg.Grow(period - avg.buffer.Len())
	} else if period < avg.buffer.Len() {
		return avg.Shrink(avg.buffer.Len() - period)
	}

	return nil
}

// Increase the period of the moving average by the given amount.
func (avg *MovingAverage[N]) Grow(amount int) {
	avg.buffer.Grow(amount)
}

// Decrease the period of the moving average by the given amount.
// If this causes the period to become less than or equal to zero
// this function returns an error.
func (avg *MovingAverage[N]) Shrink(amount int) error {
	prevSize := avg.elements

	removed, err := avg.buffer.Shrink(amount)
	if err != nil {
		return err
	}

	if avg.elements > avg.buffer.Len() {
		avg.elements = avg.buffer.Len()
	}

	avg.recalculateAvgRemoved(removed, prevSize)

	return nil
}

func (avg *MovingAverage[N]) recalculateAvgRemoved(removed []N, prevSize int) {
	for i := 0; i < len(removed); i++ {
		avg.currAvg -= removed[i] / N(prevSize)
	}

	avg.currAvg /= N(avg.elements)
	avg.currAvg *= N(prevSize)
}
