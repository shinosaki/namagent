package ndgr

import "log"

// この辺参考に
// https://github.com/rinsuki-lab/ndgr-reader/blob/main/src/protobuf-stream-reader.ts
// https://github.com/tsukumijima/NDGRClient/blob/master/ndgr_client/protobuf_stream_reader.py

type ProtobufStreamReader struct {
	buffer []byte
}

func NewProtobufStreamReader() *ProtobufStreamReader {
	return &ProtobufStreamReader{
		buffer: make([]byte, 0),
	}
}

func (r *ProtobufStreamReader) AddNewChunk(chunk []byte) {
	r.buffer = append(r.buffer, chunk...)
}

/*
reads a variable-length integer from buffer
:return offset, varInt, ok
*/
func (r *ProtobufStreamReader) readVarInt() (int, int, bool) {
	offset := 0
	result := 0
	i := 0

	for {
		if offset >= len(r.buffer) {
			return 0, 0, false
		}

		current := r.buffer[offset]
		result |= int(current&0x7F) << i
		offset++
		i += 7

		if current&0x80 == 0 {
			break
		}
	}

	return offset, result, true
}

func (r *ProtobufStreamReader) UnshiftChunk() ([]byte, bool) {
	offset, varInt, ok := r.readVarInt()
	if !ok {
		return nil, false
	}

	if offset+varInt > len(r.buffer) {
		log.Printf("needs %d bytes, but only %d bytes", offset+varInt, len(r.buffer))
		return nil, false
	}

	message := r.buffer[offset : offset+varInt]
	r.buffer = append(r.buffer[:0], r.buffer[offset+varInt:]...)
	return message, true
}
