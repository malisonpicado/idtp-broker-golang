package utils

import (
	"fmt"
	"slices"
	"testing"
)

func TestSizeOf(t *testing.T) {
	types := []string{"BOOL", "UINT_8", "UINT_16", "UINT_32", "UINT_64", "INT_8", "INT_16", "INT_32", "INT_64", "FLOAT_32", "FLOAT_64", "ANY_1", "ANY_2"}
	input := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	expect := []byte{1, 1, 2, 4, 8, 1, 2, 4, 8, 4, 8, 0, 0}

	for i, tp := range types {
		t.Run(tp, func(t *testing.T) {
			ans := SizeOf(input[i])

			if ans != expect[i] {
				t.Fatal("Required:", input[i], "| Got:", ans)
			}
		})
	}
}

func TestHasIndex(t *testing.T) {
	list := []uint32{0, 1, 2, 3, 6, 7, 10, 12, 20, 30, 100, 124, 0xFF00A0, 0xF3021}
	indexes := []uint32{33, 890, 1_000_000, 2, 6, 30, 55, 16711840, 995361}
	expect := []bool{false, false, false, true, true, true, false, true, true}

	if len(indexes) != len(expect) {
		panic("Input list and expected output list must have same length")
	}

	for i, ix := range indexes {
		t.Run(fmt.Sprint("Searching index:", ix), func(t *testing.T) {
			ans := HasIndex(ix, list)

			if ans != expect[i] {
				t.Fatal("Required:", expect[i], "| Got:", ans)
			}
		})
	}

	t.Run("Searching on nil list", func(t *testing.T) {
		ans := HasIndex(32, nil)

		if ans != false {
			t.Fatal("Required:", false, "| Got:", ans)
		}
	})

	t.Run("Searching on empty list", func(t *testing.T) {
		ans := HasIndex(32, make([]uint32, 0))

		if ans != false {
			t.Fatal("Required:", false, "| Got:", ans)
		}
	})
}

func TestCompactIndex(t *testing.T) {
	indexes := []uint32{0, 255, 256, 65_535, 65_536, 16_777_215, 16_777_216, 20_555_233, 4_294_967_295}
	expectLen := []byte{0, 0, 1, 1, 2, 2, 3, 3, 3}
	expectResult := [][]byte{
		{0x00},
		{0xFF},
		{0x01, 0x00},
		{0xFF, 0xFF},
		{0x01, 0x00, 0x00},
		{0xFF, 0xFF, 0xFF},
		{0x01, 0x00, 0x00, 0x00},
		{0x01, 0x39, 0xA5, 0xE1},
		{0xFF, 0xFF, 0xFF, 0xFF},
	}

	for i, ix := range indexes {
		t.Run(fmt.Sprint("Testing for index value:", ix), func(t *testing.T) {
			resIlen, resBytes := CompactIndex(ix)

			if resIlen != expectLen[i] {
				t.Fatal("[IndexLength] Required:", expectLen[i], "| Got:", resIlen)
			}

			if slices.Compare(expectResult[i], resBytes) != 0 {
				t.Fatal("[IndexBytes] Required:", expectResult[i], "| Got:", resBytes)
			}
		})
	}
}
