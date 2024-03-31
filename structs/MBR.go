package structs

type MBR struct {
	Tamaño     int32        // 4 bytes
	Fecha      [10]byte     // 10 bytes
	Signature  int32        // 4 bytes
	Fit        [1]byte      // 1 byte
	Partitions [4]Partition // 140 byte
} //159
