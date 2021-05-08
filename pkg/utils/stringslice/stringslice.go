package stringslice

// Contains check the slice contains str or not, if exist return true, else return false.
func Contains(slice []string, str string) bool {
	for _, value := range slice {
		if value == str {
			return true
		}
	}
	return false
}

// Delete str from slice, if success return true, else return false.
func Delete(slice *[]string, str string) bool {
	for i := 0; i < len(*slice); i++ {
		if (*slice)[i] == str {
			*slice = append((*slice)[:i], (*slice)[i+1:]...)
			return true
		}
	}
	return false
}
