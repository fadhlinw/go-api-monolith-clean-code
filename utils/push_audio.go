package utils

func ConvertPushCode(code int) string {

	describe := "pushed"
	switch code {
	case -2:
		describe = "pushed"
	case -1:
		describe = "updating"
	case 0:
		describe = "success"
	case 1:
		describe = "failed" // Hash file berbeda, file yg diterima berbeda
	case 2:
		describe = "failed" // File .tar.gz ada masalah dalam proses ekstrak
	case 3:
		describe = "failed" // Ada file amr kosong, (ada salahsatu file amr sizenya 0byte)
	case 4:
		describe = "failed" // Link url bermasalah, tidak dapat download file
	default:
		return "pushed"
	}

	return describe
}
