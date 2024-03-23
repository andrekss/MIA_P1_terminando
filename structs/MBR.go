package structs

type MBR struct {
	Tamaño     int32
	Fecha      [10]byte
	Signature  int32
	Fit        [1]byte
	Partitions [4]Partition
}
