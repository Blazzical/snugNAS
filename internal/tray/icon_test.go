package tray

import (
	"encoding/binary"
	"testing"
)

func TestSnugnasIcon(t *testing.T) {
	b := snugnasIcon()

	// Expected total: 6 (ICONDIR) + 16 (ICONDIRENTRY) + 40 (BITMAPINFOHEADER)
	// + 16*16*4 (pixels) + 16*16/8 (AND mask) = 22 + 40 + 1024 + 32 = 1118
	if len(b) != 1118 {
		t.Fatalf("icon length = %d, want 1118", len(b))
	}

	// ICONDIR
	if got := binary.LittleEndian.Uint16(b[0:]); got != 0 {
		t.Errorf("ICONDIR reserved = %d, want 0", got)
	}
	if got := binary.LittleEndian.Uint16(b[2:]); got != 1 {
		t.Errorf("ICONDIR type = %d, want 1 (icon)", got)
	}
	if got := binary.LittleEndian.Uint16(b[4:]); got != 1 {
		t.Errorf("ICONDIR count = %d, want 1", got)
	}

	// ICONDIRENTRY width/height
	if b[6] != 16 || b[7] != 16 {
		t.Errorf("icon dimensions = %dx%d, want 16x16", b[6], b[7])
	}
	if got := binary.LittleEndian.Uint16(b[12:]); got != 32 {
		t.Errorf("bit count = %d, want 32", got)
	}

	// BITMAPINFOHEADER at offset 22
	if got := binary.LittleEndian.Uint32(b[22:]); got != 40 {
		t.Errorf("BITMAPINFOHEADER size = %d, want 40", got)
	}
	if got := binary.LittleEndian.Uint32(b[26:]); got != 16 {
		t.Errorf("biWidth = %d, want 16", got)
	}
	// biHeight is doubled (image + AND mask) in ICO
	if got := binary.LittleEndian.Uint32(b[30:]); got != 32 {
		t.Errorf("biHeight = %d, want 32 (16+16 mask)", got)
	}
}
