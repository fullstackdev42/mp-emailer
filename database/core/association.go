package core

import "gorm.io/gorm"

type GormAssociation struct {
	association *gorm.Association
}

func NewAssociation(association *gorm.Association) AssociationInterface {
	return &GormAssociation{association: association}
}

func (a *GormAssociation) Find(out interface{}) error {
	return a.association.Find(out)
}

func (a *GormAssociation) Append(values ...interface{}) error {
	return a.association.Append(values...)
}

func (a *GormAssociation) Replace(values ...interface{}) error {
	return a.association.Replace(values...)
}

func (a *GormAssociation) Delete(values ...interface{}) error {
	return a.association.Delete(values...)
}

func (a *GormAssociation) Clear() error {
	return a.association.Clear()
}

func (a *GormAssociation) Count() int64 {
	return a.association.Count()
}
