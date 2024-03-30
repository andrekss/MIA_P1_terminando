package sistemaarchivos

type SuperBloque struct {
	s_filesystem_type   int32
	s_inodes_count      int32
	s_blocks_count      int32
	s_free_blocks_count int32
	s_free_inodes_count int32
	s_mtime             [10]byte
	s_umtime            [10]byte
	s_mnt_count         int32
	s_magic             int32
	s_inode_s           int32
	s_block_s           int32
	s_firts_ino         int32
	s_first_blo         int32
	s_bm_inode_start    int32
	s_bm_block_start    int32
	s_inode_start       int32
	s_block_start       int32
}
