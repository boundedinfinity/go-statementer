package util

import "errors"

type Node[T any] struct {
	Parent    *Node[T]
	Childrend []Node[T]
}

var ErrTreeWalkDone = errors.New("tree walk complete")

func (this Node[T]) Walk(node Node[T], process func(Node[T]) error) error {
	var err error

	err = process(this)
	if err != nil {
		if errors.Is(err, ErrTreeWalkDone) {
			return nil
		} else {
			return err
		}
	}

	for _, child := range this.Childrend {
		err = child.Walk(child, process)
		if err != nil {
			if errors.Is(err, ErrTreeWalkDone) {
				return nil
			} else {
				return err
			}
		}
	}

	return nil
}
