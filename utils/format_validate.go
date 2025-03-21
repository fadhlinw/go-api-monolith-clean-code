package utils

func IsMp3Format(data []byte) bool {
	if len(data) < 3 {
		return false
	}

	// Cek header ID3
	if string(data[:3]) == "ID3" {
		return true
	}

	// Cek frame header MPEG Audio (sync bits)
	return data[0] == 0xFF && (data[1]&0xE0) == 0xE0
}

func IsTarGzFormat(data []byte) bool {
	if len(data) < 4 {
		return false
	}

	// Cek header GZIP (file .gz)
	if data[0] == 0x1F && data[1] == 0x8B && data[2] == 0x08 {
		return true
	}

	// Cek header TAR (file .tar)
	if len(data) >= 512 && string(data[257:262]) == "ustar" {
		return true
	}

	return false
}
