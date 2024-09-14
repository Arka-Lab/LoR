package tools

func RandomSet(data []string, k int) ([]string, []string) {
	arr := make([]string, len(data))
	rnd := make([]int, 0)
	copy(arr, data)

	for i := 0; i < k; i++ {
		if len(rnd) == 0 {
			rnd = SHA256Arr(arr[i:])
		}
		index := rnd[0] % (len(arr) - i)
		arr[i], arr[index], rnd = arr[index], arr[i], rnd[1:]
	}
	return arr[:k], arr[k:]
}
