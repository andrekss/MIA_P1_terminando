package structs

type Partition struct { // partición 1 ó partición primaria
	Status [1]byte
	Tipo   [1]byte
	Fit    [1]byte
	Start  int32
	Size   int32
	Name   [16]byte
	Corr   int32
	Id     [4]byte
}
