package sistemaarchivos

type Inodo struct {
	I_uid   int32
	I_gid   int32
	I_size  int32
	I_atime [16]byte
	I_ctime [16]byte
	I_mtime [16]byte
	I_block [16]int32
	I_type  int32
	I_perm  int32
}
