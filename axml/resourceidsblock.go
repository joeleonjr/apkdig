package axml

/*
 * Copyright (c) 2014 Floor Terra <floort@gmail.com>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

/* +------------------------------------+
 * | Type             uint32            |
 * | Size             uint32            |
 * +------------------------------------+
 * | +--------------------------------+ |
 * | | Id             uint32          | |
 * | +--------------------------------+ |
 * |      Repeat Size/4 - 2 times       |
 * +------------------------------------+
 * |
 * +------------------------------------+
 */
type ResourceIdsBlock struct {
	AxmlBlock
	Ids []uint32
}

func (b *ResourceIdsBlock) UnmarshalBinary(data []byte) error {
	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.LittleEndian, &b.Type); err != nil {
		return err
	}
	if b.Type != CHUNK_RESOURCEIDS {
		return fmt.Errorf("Expected type=%X, got type=%X", CHUNK_RESOURCEIDS, b.Type)
	}
	if err := binary.Read(reader, binary.LittleEndian, &b.Size); err != nil {
		return err
	}
	b.Ids = make([]uint32, b.Size/4-2)
	for i := uint32(0); i < b.Size/4-2; i++ {
		if err := binary.Read(reader, binary.LittleEndian, &b.Ids[i]); err != nil {
			return err
		}
	}
	return nil
}

func (b ResourceIdsBlock) MarshalBinary() (data []byte, err error) {
	buf := bytes.NewBuffer(nil)
	if err := binary.Write(buf, binary.LittleEndian, &b.Type); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, &b.Size); err != nil {
		return nil, err
	}
	for i := range b.Ids {
		if err := binary.Write(buf, binary.LittleEndian, &b.Ids[i]); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func ReadResourceIdsBlock(reader io.ReadSeeker, size uint32, offset int64) (rid ResourceIdsBlock, err error) {
	rid.Type = CHUNK_RESOURCEIDS
	rid.Size = size
	rid.Offset = offset
	reader.Seek(offset, 0)
	rid.Ids = make([]uint32, size/4-2)
	for i := uint32(0); i < size/4-2; i++ {
		binary.Read(reader, binary.LittleEndian, &rid.Ids[i])
	}
	return rid, nil
}
