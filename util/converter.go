package util

import "strconv"

func ConvertIntToString(value int) string {
	return strconv.Itoa(value)
}

func ConvertInt32ToString(value int32) string {
	// 10 is decimal number
	return strconv.FormatInt(int64(value), 10)
}

func ConvertInt64ToString(value int64) string {
	// 10 is decimal number
	return strconv.FormatInt(value, 10)
}

func ConvertUint32ToString(value uint32) string {
	// 10 is decimal number
	return strconv.FormatInt(int64(value), 10)
}

func ConvertUint64ToString(value uint64) string {
	// 10 is decimal number
	return strconv.FormatInt(int64(value), 10)
}

func ConvertStringToInt(value string) (int, error) {
	converted, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return converted, nil
}

func ConvertStringToInt32(value string) (int32, error) {
	// 10 is decimal number, 32 is int32
	converted, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(converted), nil
}

func ConvertStringToInt64(value string) (int64, error) {
	// 10 is decimal number, 64 is int64
	converted, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return converted, nil
}
