package queue

type item interface {
	process() error
}

// type retriableItem interface {
// 	item
// 	retry() bool
// }
