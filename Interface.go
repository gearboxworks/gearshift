package barnacle

import (
	"strings"
	"os"
)

type Interface struct {
	namespace string
	name      string
}

func NewInterface(sbi string) *Interface {
	bi := new(Interface)
	bi.parse(sbi)
	return bi
}

func (bi *Interface) QualifiedName() string {
	return bi.namespace + "/" + bi.name
}

func (bi Interface) CommandsPath(leaf ...string) (path string, err error) {
	l := ""
	if len(leaf) == 1 {
		l = "/" + strings.TrimLeft(leaf[0], "/")
	}
	d, err := os.Getwd()
	return d + "/files/etc/barnacle/" + bi.QualifiedName() + l, err
}

func (bi *Interface) parse(sbi string) {
	var parts []string = strings.Split(sbi, "/")
	bi.namespace = parts[0]
	bi.name = parts[1]
}

