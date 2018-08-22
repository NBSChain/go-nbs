package application

import (
	"os"
)

type Application interface {
	AddFile(file *os.File) error
}
