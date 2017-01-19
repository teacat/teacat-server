package datastore

import (
	"github.com/TeaMeow/KitSvc/model"
	"github.com/jinzhu/gorm"
)

/*func (db *datastore) Assign(p *model.Permission) error {
	//return
}*/

func (db *datastore) Can(p *model.Permission) bool {

	action := p.Action
	p.Action = 0

	var perm model.Permission

	if db.Where(&p).First(&perm).Error == gorm.ErrRecordNotFound {
		return false
	}

	return (perm.Action & action) != 0
}
