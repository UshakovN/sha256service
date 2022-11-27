package sha256

import (
  "encoding/binary"
)

func (hash *SHA256) Sum(message []byte, secret string) [32]byte {
  hash.resetState()

  message = applySecret(message, secret)

  initial := hash.getInitial()

  h0 := initial[0]
  h1 := initial[1]
  h2 := initial[2]
  h3 := initial[3]
  h4 := initial[4]
  h5 := initial[5]
  h6 := initial[6]
  h7 := initial[7]

  k := hash.getK()

  padded := append(message, 0x80)
  if len(padded)%64 < 56 {
    suffix := make([]byte, 56-(len(padded)%64))
    padded = append(padded, suffix...)
  } else {
    suffix := make([]byte, 64+56-(len(padded)%64))
    padded = append(padded, suffix...)
  }
  msgLen := len(message) * 8
  bs := make([]byte, 8)
  binary.BigEndian.PutUint64(bs, uint64(msgLen))
  padded = append(padded, bs...)

  var broken [][]byte
  for i := 0; i < len(padded)/64; i++ {
    broken = append(broken, padded[i*64:i*64+63])
  }
  for _, chunk := range broken {
    var w []uint32
    for i := 0; i < 16; i++ {
      w = append(w, binary.BigEndian.Uint32(chunk[i*4:i*4+4]))
    }
    w = append(w, make([]uint32, 48)...)

    for i := 16; i < 64; i++ {
      s0 := rotR(w[i-15], 7) ^ rotR(w[i-15], 18) ^ (w[i-15] >> 3)
      s1 := rotR(w[i-2], 17) ^ rotR(w[i-2], 19) ^ (w[i-2] >> 10)
      w[i] = w[i-16] + s0 + w[i-7] + s1
    }

    a := h0
    b := h1
    c := h2
    d := h3
    e := h4
    f := h5
    g := h6
    h := h7

    for i := 0; i < 64; i++ {
      S1 := rotR(e, 6) ^ rotR(e, 11) ^ rotR(e, 25)
      ch := (e & f) ^ ((^e) & g)
      t1 := h + S1 + ch + k[i] + w[i]
      S0 := rotR(a, 2) ^ rotR(a, 13) ^ rotR(a, 22)
      maj := (a & b) ^ (a & c) ^ (b & c)
      t2 := S0 + maj

      h = g
      g = f
      f = e
      e = d + t1
      d = c
      c = b
      b = a
      a = t1 + t2
    }

    h0 = h0 + a
    h1 = h1 + b
    h2 = h2 + c
    h3 = h3 + d
    h4 = h4 + e
    h5 = h5 + f
    h6 = h6 + g
    h7 = h7 + h
  }

  digestBytes := [][]byte{
    toB(h0),
    toB(h1),
    toB(h2),
    toB(h3),
    toB(h4),
    toB(h5),
    toB(h6),
    toB(h7),
  }

  var digest []byte
  digestArr := [32]byte{}
  for i := 0; i < 8; i++ {
    digest = append(digest, digestBytes[i]...)
  }
  copy(digestArr[:], digest[0:32])

  return digestArr
}

func toB(i uint32) []byte {
  bs := make([]byte, 4)
  binary.BigEndian.PutUint32(bs, i)
  return bs
}

func rotR(n uint32, d uint) uint32 {
  return (n >> d) | (n << (32 - d))
}

func applySecret(message []byte, secret string) []byte {
  b := []byte(secret)
  return append(message, b...)
}