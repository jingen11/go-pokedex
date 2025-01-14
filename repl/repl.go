package repl

import "strings"

func cleanInput(text string) []string {
	sanitised := strings.Trim(text, " ")

	arr := strings.Split(sanitised, " ")

	res := []string{}

	for _, item := range arr {
		if item != "" {
			res = append(res, strings.ToLower(item))
		}
	}

	return res
}
