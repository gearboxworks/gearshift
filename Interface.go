package main

import (
	"strings"
	"os"
)

type BarnacleInterface struct {
	namespace string
	name string
}

func NewBarnacleInterface(sbi string) *BarnacleInterface {
	bi := new(BarnacleInterface)
	bi.parse(sbi)
	return bi
}

func (bi *BarnacleInterface) QualifiedName() string {
	return bi.namespace + "/" + bi.name
}

func (bi BarnacleInterface) CommandsPath(leaf ...string) (path string, err error) {
	l := ""
	if len(leaf) == 1 {
		l = "/" + strings.TrimLeft(leaf[0],"/")
	}
	d, err := os.Getwd()
	return d + "/files/etc/barnacle/" + bi.QualifiedName() + l, err
}

func (bi *BarnacleInterface) parse(sbi string) BarnacleInterface {
	var parts []string = strings.Split(sbi,"/")
	return BarnacleInterface{
		namespace: parts[0],
		name: parts[1],
	}
}

