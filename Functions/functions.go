package functions

import (
	"sync"

	u "github.com/AdityaSubrahmanyaBhat/golang/dashDB/models/User"
)

type Functions interface{
	CreateNewDB(dir string)
}

func CreateNewDB(dir string) (){

}

func Write(collection string, recordItem string, user u.User) error{

}

func ReadAll(collection string) (){

}

func Delete(collection string, recordName string)error{}

func GetMutex() *sync.Mutex{}

func CreateMutex() *sync.Mutex{}