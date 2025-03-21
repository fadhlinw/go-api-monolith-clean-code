package constants

const (
	ERROR_GETTING_USER        = "Gagal mendapatkan pengguna"
	ERROR_UPDATING_USER       = "Gagal memperbarui data pengguna"
	ERROR_DELETING_USER       = "Gagal menghapus data pengguna"
	ERROR_MARSHALING_USER     = "Gagal penyusunan data pengguna"
	ERROR_USER_NOT_FOUND      = "Pengguna tidak ditemukan"
	ERROR_USER_NOT_AUTHORIZED = "Pengguna tidak terotorisasi"

	ERROR_CREATING_OTP        = "Gagal membuat kode otp"
	ERROR_UPDATING_OTP        = "Gagal memperbarui status otp"
	ERROR_GETTING_OTP_BY_CODE = "Gagal mendapatkan otp berdasarkan kode"
	ERROR_SENDING_EMAIL       = "Gagal mengirim email"

	ERROR_NO_FILES_DOWNLOADED = "Tidak ada file yang terdownload"
	ERROR_FILE_TOO_LARGE      = "Ukuran file terlalu besar, lebih dari 500kb"

	ERROR_UPLOADING_FILE      = "Gagal mengupload file"
	ERROR_INVALID_FILE_FORMAT = "Format file tidak valid"
)
