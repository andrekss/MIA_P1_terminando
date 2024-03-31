package structs

type Partition struct { // partición 1 ó partición primaria
	Status [1]byte  // 1 byte
	Tipo   [1]byte  // 1 byte
	Fit    [1]byte  // 1 byte
	Start  int32    // 4 byte
	Size   int32    // 4 byte
	Name   [16]byte // 16 byte
	Corr   int32    // 4 byte
	Id     [4]byte  // 4 byte
} // 35
