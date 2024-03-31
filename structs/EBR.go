package structs

type EBR struct {
	part_mount [1]byte  // 1 byte
	fit        [1]byte  // 1 byte
	Start      int32    // 4 byte
	part_s     int32    // 4 byte
	next       int32    // 4 byte
	part_name  [16]byte // 16 byte
} // 30 byte
