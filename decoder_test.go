package gsp_test

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/rand/v2"
	"testing"
	"unsafe"

	"github.com/samborkent/gsp"
)

func TestDecoderDecode(t *testing.T) {
	t.Parallel()

	N := 10

	t.Run("uint8 mono", func(t *testing.T) {
		t.Parallel()

		input := make([]uint8, N)
		for i := range N {
			input[i] = uint8(rand.UintN(math.MaxUint8))
		}

		expected := make([]uint8, N)
		for i := range N {
			expected[i] = input[i]
		}

		testDecodeMono(t, input, expected)
	})

	t.Run("int8 mono", func(t *testing.T) {
		t.Parallel()

		data := make([]int8, N)
		for i := range N {
			data[i] = int8((2*rand.Float32() - 1) * math.MinInt8)
		}

		input := make([]byte, N)
		expected := make([]uint8, N)
		for i := range N {
			input[i] = byte(data[i])
			expected[i] = input[i]
		}

		testDecodeMono(t, input, expected)
	})

	t.Run("uint16 mono", func(t *testing.T) {
		t.Parallel()

		data := make([]uint16, N)
		for i := range N {
			data[i] = uint16(rand.UintN(math.MaxUint16))
		}

		input := make([]byte, 2*N)
		expected := make([]uint16, N)
		for i := range N {
			binary.LittleEndian.PutUint16(input[2*i:2*i+2], data[i])
			expected[i] = data[i]
		}

		testDecodeMono(t, input, expected)
	})

	t.Run("int16 mono", func(t *testing.T) {
		t.Parallel()

		data := make([]int16, N)
		for i := range N {
			data[i] = int16((2*rand.Float32() - 1) * math.MinInt16)
		}

		input := make([]byte, 2*N)
		expected := make([]int16, N)
		for i := range N {
			binary.LittleEndian.PutUint16(input[2*i:2*i+2], uint16(data[i]))
			expected[i] = data[i]
		}

		testDecodeMono(t, input, expected)
	})

	t.Run("uint32 mono", func(t *testing.T) {
		t.Parallel()

		data := make([]uint32, N)
		for i := range N {
			data[i] = rand.Uint32()
		}

		input := make([]byte, 4*N)
		expected := make([]uint32, N)
		for i := range N {
			binary.LittleEndian.PutUint32(input[4*i:4*i+4], data[i])
			expected[i] = data[i]
		}

		testDecodeMono(t, input, expected)
	})

	t.Run("int32 mono", func(t *testing.T) {
		t.Parallel()

		data := make([]int32, N)
		for i := range N {
			data[i] = int32((2*rand.Float32() - 1) * math.MinInt32)
		}

		input := make([]byte, 4*N)
		expected := make([]int32, N)
		for i := range N {
			binary.LittleEndian.PutUint32(input[4*i:4*i+4], uint32(data[i]))
			expected[i] = data[i]
		}

		testDecodeMono(t, input, expected)
	})

	t.Run("float32 mono", func(t *testing.T) {
		t.Parallel()

		data := make([]float32, N)
		for i := range N {
			data[i] = 2*rand.Float32() - 1
		}

		input := make([]byte, 4*N)
		expected := make([]float32, N)
		for i := range N {
			binary.LittleEndian.PutUint32(input[4*i:4*i+4], math.Float32bits(data[i]))
			expected[i] = data[i]
		}

		testDecodeMono(t, input, expected)
	})

	t.Run("float64 mono", func(t *testing.T) {
		t.Parallel()

		data := make([]float64, N)
		for i := range N {
			data[i] = 2*rand.Float64() - 1
		}

		input := make([]byte, 8*N)
		expected := make([]float64, N)
		for i := range N {
			binary.LittleEndian.PutUint64(input[8*i:8*i+8], math.Float64bits(data[i]))
			expected[i] = data[i]
		}

		testDecodeMono(t, input, expected)
	})
}

func testDecodeMono[T gsp.Type](t *testing.T, input []byte, want []T) {
	t.Helper()

	decoder := gsp.NewDecoder[T, T](bytes.NewReader(input))

	if decoder.Channels() != 1 {
		t.Errorf("wrong number of channels: got '%d', want '%d'", decoder.Channels(), 1)
	}

	byteSize := int(unsafe.Sizeof(T(0)))

	if decoder.ByteSize() != byteSize {
		t.Errorf("wrong byte size: got '%d', want '%d'", decoder.ByteSize(), byteSize)
	}

	samples := make([]T, len(input)/byteSize)

	err := decoder.Decode(samples)
	if err != nil {
		t.Fatalf("decoding samples: error: %s", err.Error())
	}

	if len(samples) != len(want) {
		t.Fatalf("missing samples: got '%d', want '%d'", len(samples), len(want))
	}

	for i := range samples {
		if samples[i] != want[i] {
			t.Errorf("sample mismatch at index '%d': got '%v', want '%v'", i, samples[i], want[i])
		}
	}
}
