package entityUtils

import (
	"github.com/astaxie/beego/orm"
)

func DefaultRegisterModel(entityInstance interface{}) {
	orm.RegisterModel(entityInstance)
}
