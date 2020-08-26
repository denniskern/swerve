package config

func hasDuplicates(arr []int) bool {
	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr) && j != i; j++ {
			if arr[i] == arr[j] {
				return true
			}
		}
	}
	return false
}
