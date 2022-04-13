package avformat

import (
	"reflect"
	"unsafe"
)

func (ctx *AVFormatContext) NbStreams() uint {
	return (uint)(ctx.nb_streams)
}

func (ctx *AVFormatContext) Streams() []*AVStream {
	header := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(ctx.streams)),
		Len:  int(ctx.nb_streams),
		Cap:  int(ctx.nb_streams),
	}
	return *(*[]*AVStream)(unsafe.Pointer(&header))
}
