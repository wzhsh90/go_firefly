package utils

func SqlLike(name string) string {
	var realName = ""
	if name != "" {
		realName = "%" + name + "%"
	}
	return realName
}
