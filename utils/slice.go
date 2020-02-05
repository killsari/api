package utils

func InStrSlice(item string, arr *[]string) (in bool) {
	for _, v := range *arr {
		if item == v {
			return true
		}
	}
	return
}

func InArray(item interface{}, arr *[]interface{}) (in bool) {
	for _, v := range *arr {
		if item == v {
			return true
		}
	}
	return
}

func Map2KeyArr(m map[uint64]string) (keyArr []uint64) {
	for key, _ := range m {
		keyArr = append(keyArr, key)
	}
	return
}

func Map2ValArr(m map[uint64]string) (valArr []string) {
	for _, val := range m {
		valArr = append(valArr, val)
	}
	return
}

func MapInvert(m map[uint64]string) (inventM map[string]uint64) {
	inventM = map[string]uint64{}
	for key, val := range m {
		inventM[val] = key
	}
	return
}
