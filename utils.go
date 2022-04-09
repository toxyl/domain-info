package main

import "fmt"

func uniqueStrings(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func printStatusLn(format string, a ...interface{}) {
	fmt.Printf("\033[s\033[2K"+format+"\033[u", a...)
}
