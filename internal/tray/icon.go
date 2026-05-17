package tray

import "encoding/binary"

// snugnasIcon builds a 16x16 32-bit ICO of a solid square in the snugNAS
// accent colour (#4ea1ff). Hand-rolled to avoid pulling in an ICO encoder.
func snugnasIcon() []byte {
	const w, h = 16, 16
	const headerOff = 22 // ICONDIR (6) + ICONDIRENTRY (16)
	const bmpHdrOff = headerOff
	const pixelsOff = bmpHdrOff + 40

	pixelBytes := w * h * 4  // 32bpp BGRA
	maskBytes := (w * h) / 8 // 1bpp AND mask
	bmpBytes := 40 + pixelBytes + maskBytes
	total := headerOff + bmpBytes

	b := make([]byte, total)

	// ICONDIR
	binary.LittleEndian.PutUint16(b[0:], 0) // reserved
	binary.LittleEndian.PutUint16(b[2:], 1) // type = icon
	binary.LittleEndian.PutUint16(b[4:], 1) // count

	// ICONDIRENTRY
	b[6] = w
	b[7] = h
	b[8] = 0                                                // color count
	b[9] = 0                                                // reserved
	binary.LittleEndian.PutUint16(b[10:], 1)                // planes
	binary.LittleEndian.PutUint16(b[12:], 32)               // bit count
	binary.LittleEndian.PutUint32(b[14:], uint32(bmpBytes)) // bytes in res
	binary.LittleEndian.PutUint32(b[18:], headerOff)        // image offset

	// BITMAPINFOHEADER
	bh := b[bmpHdrOff : bmpHdrOff+40]
	binary.LittleEndian.PutUint32(bh[0:], 40)          // biSize
	binary.LittleEndian.PutUint32(bh[4:], uint32(w))   // biWidth
	binary.LittleEndian.PutUint32(bh[8:], uint32(h*2)) // biHeight (image + mask)
	binary.LittleEndian.PutUint16(bh[12:], 1)          // biPlanes
	binary.LittleEndian.PutUint16(bh[14:], 32)         // biBitCount
	// remaining fields stay zero (BI_RGB, no compression)

	// Pixel data: BGRA, bottom-up. #4ea1ff -> R=78, G=161, B=255, A=255.
	pixels := b[pixelsOff : pixelsOff+pixelBytes]
	for i := 0; i < w*h; i++ {
		pixels[i*4+0] = 255 // B
		pixels[i*4+1] = 161 // G
		pixels[i*4+2] = 78  // R
		pixels[i*4+3] = 255 // A
	}
	// AND mask is all-zero (opaque), which is what `make` already gave us.

	return b
}
