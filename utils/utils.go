package utils

// Base62Encode encodes a given integer into a base62 string.
// It uses the characters 0-9, a-z, and A-Z for encoding.
//
// Parameters:
//
//	num (int64): The integer to be encoded.
//
// Returns:
//
//	string: The base62 encoded string.
func Base62Encode(num int64) string {
	if num == 0 {
		return "0"
	}
	if num < 0 {
		return ""
	}
	const base62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := ""
	for num > 0 {
		result = string(base62[num%62]) + result
		num /= 62
	}
	return result
}
