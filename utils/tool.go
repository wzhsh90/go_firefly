package utils

import "strconv"

func SqlLike(name string) string {
	var realName = ""
	if name != "" {
		realName = "%" + name + "%"
	}
	return realName
}
func ParseUnInt(val string) uint {
	pval, _ := strconv.ParseUint(val, 10, 64)
	return uint(pval)
}
func ParseInt(val string) int {
	pval, _ := strconv.ParseInt(val, 10, 64)
	return int(pval)
}
