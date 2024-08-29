package buffer

type Buffer  struct 
{
	Buf  []byte
	Len int
}




func (b *Buffer) BAppend (data []byte) () {
	b.Buf = append(b.Buf, data...)
	b.Len += len(data)
 
}


func (b* Buffer) BFree () {
	b.Buf = nil 
}



func  InitBuffer () *Buffer  {
	return  &Buffer{
		Buf: nil,
		Len: 0,
	}
}