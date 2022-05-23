package callbacks

import "github.com/driver005/database"

type BeforeCreateInterface interface {
	BeforeCreate(*database.DB) error
}

type AfterCreateInterface interface {
	AfterCreate(*database.DB) error
}

type BeforeUpdateInterface interface {
	BeforeUpdate(*database.DB) error
}

type AfterUpdateInterface interface {
	AfterUpdate(*database.DB) error
}

type BeforeSaveInterface interface {
	BeforeSave(*database.DB) error
}

type AfterSaveInterface interface {
	AfterSave(*database.DB) error
}

type BeforeDeleteInterface interface {
	BeforeDelete(*database.DB) error
}

type AfterDeleteInterface interface {
	AfterDelete(*database.DB) error
}

type AfterFindInterface interface {
	AfterFind(*database.DB) error
}
