package gocache

//ByteView ByteView
type ByteView struct {
	b []byte
}

//Len Len
func (v ByteView) Len() int {
	return len(v.b)
}

//ByteSlice return copy of byte
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
