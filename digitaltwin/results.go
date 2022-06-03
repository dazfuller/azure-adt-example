package digitaltwin

import "azure-adt-example/digitaltwin/models"

// TwinResult2 defines a result set consisting of two models.IModel types.
type TwinResult2[T1, T2 models.IModel] struct {
	Twin1 T1
	Twin2 T2
}

// NewTwinResult2 creates a new instance of the TwinResult2 type.
func NewTwinResult2[T1, T2 models.IModel](t1 *T1, t2 *T2) TwinResult2[T1, T2] {
	return TwinResult2[T1, T2]{*t1, *t2}
}

// TwinResult3 defines a result set consisting of three models.IModel types.
type TwinResult3[T1, T2, T3 models.IModel] struct {
	Twin1 T1
	Twin2 T2
	Twin3 T3
}

// NewTwinResult3 creates a new instance of the TwinResult3 type.
func NewTwinResult3[T1, T2, T3 models.IModel](t1 *T1, t2 *T2, t3 *T3) TwinResult3[T1, T2, T3] {
	return TwinResult3[T1, T2, T3]{*t1, *t2, *t3}
}
