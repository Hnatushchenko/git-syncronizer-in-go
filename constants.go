package main

type constants struct {
}

func GetDirectoryNamesToIgnore() []string {
	return []string{
		"node_modules",
		".angular",
		".git",
		"obj",
		"bin",
		".nx",
	}
}

func GetFileNamesToIgnore() []string {
	return []string{
		".gitmodules",
	}
}
