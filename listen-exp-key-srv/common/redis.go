package common

import "fmt"

const (
	//redis key
	ListenExpKeyLock = "ListenExpKeyLock"
)

func GetListenExpKeyLockKey(key string) string {
	return fmt.Sprintf("%s_%s", ListenExpKeyLock, key)
}
