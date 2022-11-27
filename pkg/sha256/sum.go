package sha256

import (
	"bytes"
	"encoding/binary"
	"math/bits"
)

const (
	sumSize       = 32
	chunkSize     = 64
	chunkSizeBits = chunkSize * 8
)

func (h *SHA256) Sum(data []byte) [sumSize]byte {
	h.resetState()

	padded := h.padMessage(data)
	chunks := h.chunks(padded)

	initial := h.initial
	k := h.k

	for _, chunk := range chunks {
		ms := h.messageSchedule(chunk)

		for i := 16; i < 64; i++ {
			s0 := bits.RotateLeft32(ms[i-15], -7) ^ bits.RotateLeft32(ms[i-15], -18) ^ (ms[i-15] >> 3)
			s1 := bits.RotateLeft32(ms[i-2], -17) ^ bits.RotateLeft32(ms[i-2], -19) ^ (ms[i-2] >> 10)
			ms[i] = ms[i-16] + s0 + ms[i-7] + s1
		}

		a := initial[0]
		b := initial[1]
		c := initial[2]
		d := initial[3]
		e := initial[4]
		f := initial[5]
		g := initial[6]
		h := initial[7]

		for j := 0; j < 64; j++ {
			S1 := bits.RotateLeft32(e, -6) ^ bits.RotateLeft32(e, -11) ^ bits.RotateLeft32(e, -25)
			ch := (e & f) ^ ((^e) & g)
			temp1 := h + S1 + ch + k[j] + ms[j]
			S0 := bits.RotateLeft32(a, -2) ^ bits.RotateLeft32(a, -13) ^ bits.RotateLeft32(a, -22)
			maj := (a & b) ^ (a & c) ^ (b & c)
			temp2 := S0 + maj
			h = g
			g = f
			f = e
			e = d + temp1
			d = c
			c = b
			b = a
			a = temp1 + temp2
		}
		initial[0] += a
		initial[1] += b
		initial[2] += c
		initial[3] += d
		initial[4] += e
		initial[5] += f
		initial[6] += g
		initial[7] += h
	}

	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, initial)

	var out [sumSize]byte
	copy(out[:], buf.Bytes())

	h.resetState()

	return out
}

func (h *SHA256) numPadZero(L int) int {
	lenInBits := L * 8
	m := lenInBits + 1 + 64
	return chunkSizeBits - m%chunkSizeBits
}

func (h *SHA256) padMessage(data []uint8) []uint8 {
	b := data
	dataLen := len(data)
	zerosLen := ((h.numPadZero(dataLen) + 1) / 8) - 1
	firstByteAfterMessage := uint8(0b10000000)

	b = append(b, firstByteAfterMessage)
	for i := 0; i < zerosLen; i++ {
		b = append(b, 0b00000000)
	}

	lenAsBytes := make([]uint8, 8)
	binary.BigEndian.PutUint64(lenAsBytes, uint64(dataLen*8))
	b = append(b, lenAsBytes...)

	return b
}

func (h *SHA256) chunks(data []uint8) [][]uint8 {
	chunks := make([][]uint8, 0)
	for i := 0; i < len(data); i += chunkSize {
		chunks = append(chunks, data[i:i+chunkSize])
	}
	return chunks
}

func (h *SHA256) messageSchedule(chunk []uint8) []uint32 {
	messageSchedule := make([]uint32, 0)
	for i := 0; i < len(chunk); i += 4 {
		messageSchedule = append(messageSchedule, binary.BigEndian.Uint32(chunk[i:i+4]))
	}
	for j := 16; j < 64; j++ {
		messageSchedule = append(messageSchedule, 0)
	}
	return messageSchedule
}
