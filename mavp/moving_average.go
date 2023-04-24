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

func New[N generics.Float](period int) MovingAverage[N] {
	return MovingAverage[N]{
		buffer:   rinbuf.New[N](period),
		elements: 0,
		currAvg:  0,
	}
}

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

func (avg *MovingAverage[N]) Current() N {
	return avg.currAvg
}

func (avg *MovingAverage[N]) Resize(period int) error {
	if period > avg.buffer.Len() {
		avg.Grow(period - avg.buffer.Len())
	} else if period < avg.buffer.Len() {
		return avg.Shrink(avg.buffer.Len() - period)
	}

	return nil
}

func (avg *MovingAverage[N]) Grow(amount int) {
	avg.buffer.Grow(amount)
}

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
