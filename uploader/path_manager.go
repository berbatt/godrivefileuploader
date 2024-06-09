package uploader

import "strings"

type Manager struct {
	pathToParentIDMap map[string]string
}

func NewManager() Manager {
	return Manager{pathToParentIDMap: make(map[string]string)}
}

func (m *Manager) SetParentID(path, parentFolderID string) {
	parent := GetFileOrFolderNameFromPath(path)
	m.pathToParentIDMap[parent] = parentFolderID
}

func (m *Manager) GetParentFolderID(path string) string {
	parent := getParentFolderName(path)
	return m.pathToParentIDMap[parent]
}

func getParentFolderName(path string) (root string) {
	paths := strings.Split(path, "/")
	if len(paths) < 2 {
		root = path
	} else {
		root = paths[len(paths)-2]
	}
	return root
}

func GetFileOrFolderNameFromPath(path string) (fileName string) {
	paths := strings.Split(path, "/")
	if len(paths) > 0 {
		return paths[len(paths)-1]
	}
	return path
}
