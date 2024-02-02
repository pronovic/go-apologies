package equality

// EqualsByValue identifies an interface that supports by-value equality via the Equals method
type EqualsByValue[T any] interface {
	// Equals Checks for value equality on the interface
	Equals(other T) bool
}

// ByValueEquals check whether two interfaces that implement EqualsByValue are equal
func ByValueEquals[T any](left EqualsByValue[T], right EqualsByValue[T]) bool {
	if left == nil && right == nil {
		return true
	} else if (left == nil && right != nil) || (left != nil && right == nil) {
		return false
	} else {
		return left.Equals(right.(T))
	}
}

// IntPointerEquals checks whether two integer pointers have the same value
func IntPointerEquals(left *int, right *int) bool {
	if left == nil && right == nil {
		return true
	} else if (left == nil && right != nil) || (left != nil && right == nil) {
		return false
	} else {
		return *left == *right
	}
}

// StringPointerEquals checks whether two string pointers have the same value
func StringPointerEquals(left *string, right *string) bool {
	if left == nil && right == nil {
		return true
	} else if (left == nil && right != nil) || (left != nil && right == nil) {
		return false
	} else {
		return *left == *right
	}
}