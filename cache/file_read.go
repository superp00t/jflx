package cache

import (
	"bytes"
	"crypto/sha512"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func index_to_real_position(index, file_size int64) (real_position int64, err error) {
	real_position = index * file_part_size
	if real_position >= file_size {
		err = fmt.Errorf("index %d is out of bounds of file size %d", index, file_size)
	}
	return
}

func index_to_part_length(index, file_size int64) (part_length int64, err error) {
	real_position := (index * file_part_size)
	if real_position >= file_size {
		err = fmt.Errorf("index %d is out of bounds of file size %d", index, file_size)
		return
	}
	if real_position+file_part_size < file_size {
		part_length = file_part_size
	} else {
		part_length = file_size - real_position
	}

	return
}

// Attempt to read a cached segment of the file. writes to the cache if not already in the cache
func (file *cache_file) read_cache_part(index int64) (data []byte, err error) {
	// First, attempt to read from cache.
	subfile_path := file.subfile_path(fmt.Sprintf("%s-%08x.part", file.hash, index))
	if _, err = os.Stat(subfile_path); err == nil {
		// if cached part exists, read from it
		var data_chunk []byte
		data_chunk, err = os.ReadFile(subfile_path)
		if err == nil {
			// if cached part can be read, verify the part's checksum
			checksum := data_chunk[0:64]
			actual_data := data_chunk[64:]
			actual_checksum := sha512.Sum512(actual_data)
			// return actual part data if checksum is valid
			if bytes.Equal(checksum, actual_checksum[:]) {
				// log.Println("Reading cached part", file.realpath, subfile_path)
				return actual_data, nil
			} else {
				log.Println("Cached file was invalid", file.realpath, subfile_path)
				// otherwise, remove the invalid part data
				os.Remove(subfile_path)
			}
		}
	} else {
		// log.Println("cached part does not exist for", file.realpath, subfile_path)
	}

	// if we can't get the part from the cache, read it from the source file and commit it to the cache.
	// read full segment/part from the source file
	var real_position int64
	var part_length int64
	real_position, err = index_to_real_position(index, file.meta.Info.InfoSize)
	if err != nil {
		err = fmt.Errorf("index miscalculation in %s: %w", file.realpath, err)
		return
	}
	part_length, err = index_to_part_length(index, file.meta.Info.InfoSize)
	if err != nil {
		err = fmt.Errorf("part length miscalculation: in %s: %w", file.realpath, err)
		return
	}
	var real_file http.File
	real_file, err = file.server.source.Open(file.meta.Realpath)
	if err != nil {
		return
	}
	if _, err = real_file.Seek(real_position, io.SeekStart); err != nil {
		return
	}
	part_data := make([]byte, part_length)
	_, err = io.ReadFull(real_file, part_data[:])
	if err != nil {
		return
	}
	real_file.Close()

	// Commit segment to cache (firstly, as a temporary file so a user can't access in the middle of writing)
	temp_path := subfile_path + ".tmp"
	part_checksum := sha512.Sum512(part_data)
	var part_file *os.File
	part_file, err = os.OpenFile(temp_path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0700)
	if err != nil {
		return
	}
	if _, err = part_file.Write(part_checksum[:]); err != nil {
		return
	}
	if _, err = part_file.Write(part_data[:]); err != nil {
		return
	}
	part_file.Close()

	// Add byte usage
	file.server.add_used_bytes(64 + int64(len(part_data)))

	// Make cached file part available to user
	if err = os.Rename(temp_path, subfile_path); err != nil {
		return
	}

	data = part_data
	err = nil

	return
}

// Attempt to read from the file cache.
func (file *cache_file) Read(data []byte) (n int, err error) {
	begin_pointer := file.pointer
	end_pointer := file.pointer + int64(len(data))

	if begin_pointer >= file.meta.Info.Size() {
		// we've reached the end of the file
		err = io.EOF
		return
	}

	// Don't read past the end of the file
	if end_pointer > file.meta.Info.Size() {
		end_pointer = file.meta.Info.Size()
	}

	begin_part_index := begin_pointer / file_part_size
	end_part_index := end_pointer / file_part_size

	part_relative_begin_offset := begin_pointer % file_part_size
	// part_relative_end_offset := end_pointer % file_part_size

	var part_data []byte

	// same part (easy!)
	if begin_part_index == end_part_index {
		part_data, err = file.read_cache_part(begin_part_index)
		if err != nil {
			return
		}

		n = copy(data, part_data[part_relative_begin_offset:])
		file.pointer += int64(n)

		return
	}

	// several parts (boy this slice is gigantic 😓)
	for part_index := begin_part_index; part_index <= end_part_index; part_index++ {
		if n >= len(data) {
			break
		}

		part_data, err = file.read_cache_part(part_index)
		if err != nil {
			break
		}

		part_begin_offset := int64(0)

		if part_index == begin_part_index {
			part_begin_offset = part_relative_begin_offset
		}

		n += copy(data[n:], part_data[part_begin_offset:])
	}

	file.pointer += int64(n)

	return
}
