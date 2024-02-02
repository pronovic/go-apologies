package equality

type EqualsByValue[T any] interface {
	// Equals Checks for value equality on the interface
	Equals(other T) bool
}

func ByValueEquals[T any](left EqualsByValue[T], right EqualsByValue[T]) bool {
	if left == nil && right == nil {
		return true
	} else if (left == nil && right != nil) || (left != nil && right == nil) {
		return false
	} else {
		return left.Equals(right.(T))
	}
}

func IntPointerEquals(left *int, right *int) bool {
	if left == nil && right == nil {
		return true
	} else if (left == nil && right != nil) || (left != nil && right == nil) {
		return false
	} else {
		return *left == *right
	}
}

func StringPointerEquals(left *string, right *string) bool {
	if left == nil && right == nil {
		return true
	} else if (left == nil && right != nil) || (left != nil && right == nil) {
		return false
	} else {
		return *left == *right
	}
}