package render

import "github.com/hyhecor/qq/opop/models"

// func first[T any](aa ...T) (b T) {
// 	for _, v := range aa {
// 		return v
// 	}
// 	return
// }

// func String[T ~string](v *T) string {
// 	if v == nil {
// 		return ""
// 	}
// 	return (string)(*v)
// }

func _map[A any, B any](fn func(a A) B) func(aa []A) []B {
	return func(aa []A) (bb []B) {
		bb = make([]B, len(aa))

		for i := range aa {
			bb[i] = fn(aa[i])
		}
		return bb
	}
}

func _ColumnToInterface(a models.Column) models.ColumnIndentifier {
	return a
}
