package gsp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"sync/atomic"
	"unsafe"
)

var ErrEncoderNotinitialized = errors.New("gsp: Encoder: not initialized")

// Encoder is a linear PCM and floating-point encoder.
// This encoder can convert samples into their binary representation.
type Encoder[F Frame[T], T Type] struct {
	w                            io.Writer
	pool                         *ByteBufferPool
	samplesEncoded, bytesWritten atomic.Int64
	channels, byteSize           int
	bigEndian                    bool
	initialized                  bool
}

func NewEncoder[F Frame[T], T Type](w io.Writer, opts ...EncodingOption) *Encoder[F, T] {
	// Apply encoding options.
	var cfg EncodingConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	var channels int

	switch frameType := any(*new(F)).(type) {
	case T:
		channels = 1
	case [2]T, Stereo[T]:
		channels = 2
	case []T:
		channels = len(frameType)
	case MultiChannel[T]:
		channels = len(frameType)
	default:
		panic("gsp: NewEncoder: unknown audio frame type")
	}

	return &Encoder[F, T]{
		w:           w,
		pool:        NewByteBufferPool(),
		channels:    channels,
		byteSize:    int(unsafe.Sizeof(T(0))),
		bigEndian:   cfg.BigEndian,
		initialized: true,
	}
}

func (e *Encoder[F, T]) Channels() int {
	return e.channels
}

func (e *Encoder[F, T]) ByteSize() int {
	return e.byteSize
}

// Encode reads from the sample slice and decodes into the internal [io.Writer].
func (e *Encoder[F, T]) Encode(s []F) error {
	if !e.initialized {
		return ErrDecoderNotinitialized
	}

	buf := e.pool.Get()
	if buf == nil {
		// Allocate in-case the pool does not have a pre-allocated entry ready.
		buf = bytes.NewBuffer(make([]byte, len(s)*e.channels*e.byteSize))
	}

	switch e.channels {
	case 1:
		// There is only a single channel, so we can safely perform this unsafe type-casting.
		samplesEncoded := e.convertMono(buf, unsafe.Slice((*T)(unsafe.Pointer(&s[0])), len(s)))
		e.samplesEncoded.Add(int64(samplesEncoded))
	case 2:
		panic("gsp: Encoder.Encode: stereo encoding not implemented")
	default:
		panic("gsp: Encoder.Encode: multi-channel encoding not implemented")
	}

	bytesWritten, err := buf.WriteTo(e.w)
	if err != nil {
		return fmt.Errorf("gsp: Encoder.Encode: writing bytes: %w", err)
	}

	e.bytesWritten.Add(int64(bytesWritten))
	e.pool.Put(buf)

	return nil
}

// convertMono writes the encoded sampled to the internal [io.Writer].
func (e *Encoder[F, T]) convertMono(buf *bytes.Buffer, src []T) int {
	switch e.byteSize {
	case 1: // 8 bit
		// Abuse overflow rules to deduce specific type.
		if isUnsigned[T]() {
			// uint8
			n, _ := buf.Write(unsafe.Slice((*uint8)(unsafe.Pointer(&src[0])), len(src)))
			return n
		}

		// int8
		for i := range len(src) {
			_ = buf.WriteByte(byte(src[i]))
		}

		return len(src)
	case 2: // 16 bit
		// uint16 & int16
		if e.bigEndian {
			for i := range len(src) {
				var data [2]byte
				binary.BigEndian.PutUint16(data[:], uint16(src[i]))
				_, _ = buf.Write(data[:])
			}
		} else {
			for i := range len(src) {
				var data [2]byte
				binary.LittleEndian.PutUint16(data[:], uint16(src[i]))
				_, _ = buf.Write(data[:])
			}
		}

		return len(src)
	case 4: // 32 bit
		// Abuse overflow rules to deduce specific type.
		if isUnsigned[T]() || isSigned[T]() {
			// uint32 & int32
			if e.bigEndian {
				for i := range len(src) {
					var data [4]byte
					binary.BigEndian.PutUint32(data[:], uint32(src[i]))
					_, _ = buf.Write(data[:])
				}
			} else {
				for i := range len(src) {
					var data [4]byte
					binary.LittleEndian.PutUint32(data[:], uint32(src[i]))
					_, _ = buf.Write(data[:])
				}
			}

			return len(src)
		}

		// float32
		if e.bigEndian {
			for i := range len(src) {
				var data [4]byte
				binary.BigEndian.PutUint32(data[:], math.Float32bits(float32(src[i])))
				_, _ = buf.Write(data[:])
			}
		} else {
			for i := range len(src) {
				var data [4]byte
				binary.LittleEndian.PutUint32(data[:], math.Float32bits(float32(src[i])))
				_, _ = buf.Write(data[:])
			}
		}

		return len(src)
	case 8: // 64 bit
		// Abuse overflow rules to deduce specific type.
		if isUnsigned[T]() || isSigned[T]() {
			// uint64 & int64
			if e.bigEndian {
				for i := range len(src) {
					var data [8]byte
					binary.BigEndian.PutUint64(data[:], uint64(src[i]))
					_, _ = buf.Write(data[:])
				}
			} else {
				for i := range len(src) {
					var data [8]byte
					binary.LittleEndian.PutUint64(data[:], uint64(src[i]))
					_, _ = buf.Write(data[:])
				}
			}

			return len(src)
		}

		// float64
		if e.bigEndian {
			for i := range len(src) {
				var data [8]byte
				binary.BigEndian.PutUint64(data[:], math.Float64bits(float64(src[i])))
				_, _ = buf.Write(data[:])
			}
		} else {
			for i := range len(src) {
				var data [8]byte
				binary.LittleEndian.PutUint64(data[:], math.Float64bits(float64(src[i])))
				_, _ = buf.Write(data[:])
			}
		}

		return len(src)
	default:
		panic("gsp: Encoder.convertMono: unknown bit size encountered")
	}
}
