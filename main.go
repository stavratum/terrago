package main

import (
	"fmt"
	"io"
	"os"

	"github.com/stavratum/terrago/filetype"
)

// for World.Flags
const WorldIsFavorite uint64 = 1 << 58

type World struct {
	Version      uint32
	HeaderFormat string

	FileRevision uint32
	Flags        uint64

	TileFrameImportant []bool

	Title string
	Seed  string
}

func (wld *World) String() string {
	return fmt.Sprintf(""+
		"terrago.World {\n"+
		"\tVersion: %d\n"+
		"\tHeaderFormat: %s\n"+
		"\n"+
		"\tTitle: %s\n"+
		"\tSeed: %s\n"+
		""+
		"}",

		wld.Version,
		wld.HeaderFormat,
		wld.Title,
		wld.Seed,
	)
}

func (wld *World) Parse(r WorldReader) (err error) {
	r.Uint32(&wld.Version)

	if wld.Version < 87 {
		return fmt.Errorf("invalid world version: %d", wld.Version)
	}
	if wld.Version >= 140 {
		switch r.String(&wld.HeaderFormat, 7); wld.HeaderFormat {
		case "relogic", "xindong":
		default:
			return fmt.Errorf("invalid world header: %s", wld.HeaderFormat)
		}

		if b := r.Byte(); b != filetype.World {
			return fmt.Errorf("invalid file type: %s", filetype.String(b))
		}

		r.Uint32(&wld.FileRevision)
		r.Uint64(&wld.Flags)
	}

	sectionCount := int16(0)
	r.Int16(&sectionCount)
	fmt.Println("sectionCount: ", sectionCount)

	sectionPointers := make([]int32, sectionCount)

	for i := range sectionCount {
		r.Int32(&sectionPointers[i])
		fmt.Println("sectionPointers[", i, "]:", sectionPointers[i])
	}

	/* read bit array */
	{
		length := int16(0)
		r.Int16(&length)

		wld.TileFrameImportant = make([]bool, length)

		var (
			data    byte
			bitMask byte = 128
		)

		for i := range length {
			// If we read the last bit mask (B1000000 = 0x80 = 128), read the next byte from the stream and start the mask over.
			// Otherwise, keep incrementing the mask to get the next bit.
			if bitMask != 128 {
				bitMask = (byte)(bitMask << 1)
			} else {
				data = r.Byte()
				bitMask = 1
			}

			// Check the mask, if it is set then set the current boolean to true
			if (data & bitMask) == bitMask {
				wld.TileFrameImportant[i] = true
			}
		}
	}

	if offset, err := r.Seek(0, io.SeekCurrent); err != nil {
		return err
	} else {
		if offset != int64(sectionPointers[0]) {
			return fmt.Errorf("unexpected position: invalid file format section")
		}
	}

	r.String(&wld.Title,
		int(r.Byte()),
	)

	r.String(&wld.Seed,
		int(r.Byte()),
	)

	return
}

func main() {
	wld := new(World)

	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	err = wld.Parse(NewReader(f))
	if err != nil {
		panic(err)
	}

	fmt.Println(wld.String())
	fmt.Println("IsFavorite:", wld.Flags&WorldIsFavorite)
}
