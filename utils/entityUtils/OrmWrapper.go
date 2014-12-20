package entityUtils

/// To use a different 'repo' just swop out all the below 'beego' and 'orm' related stuff with the other repo
/// When we switch it will not compile because in the beegoOrmRepo we make use of *OrmWrapper which
/// now wraps a different type of OrmInstance
/// Because we will not use the beego repo we could actually just comment out its code, but the better
/// solution is to just uncomment the 'beegoOrmWrapper' and replace all the OrmWrapper usages with it

import (
	"github.com/astaxie/beego/orm"
)

const (
	DEFAULT_QUERY_LIMIT = 1500
)

type OrmWrapper struct {
	OrmInstance               orm.Ormer
	LocallyManageTransactions bool
}

func CreateNewOrmWrapper(possibleParentTransactionalOrmWrapper *OrmWrapper) *OrmWrapper {
	if possibleParentTransactionalOrmWrapper == nil || possibleParentTransactionalOrmWrapper.OrmInstance == nil {
		return &OrmWrapper{
			OrmInstance:               orm.NewOrm(),
			LocallyManageTransactions: true, //True because we do not have a 'parent' ORM so we can manage the transactions ourselves
		}
	}
	return &OrmWrapper{
		//Share the OrmInstance otherwise if we use a different one they will not be part of the same DB transaction
		OrmInstance: possibleParentTransactionalOrmWrapper.OrmInstance,
		//False because we have a 'parent' ORM so we do not do transactions, instead just return error and
		//assume our 'parent' will rollback the transaction on errors
		LocallyManageTransactions: false,
	}
}

func (this *OrmWrapper) BeginTransaction() {
	//TODO: We can probably at least log an error if it happens inside one of these methods (BeginTransaction, RollbackTransaction, CommitTransaction)?
	if !this.LocallyManageTransactions {
		return
	}
	err := this.OrmInstance.Begin()
	if err != nil {
		panic(err)
	}
}

func (this *OrmWrapper) RollbackTransaction() {
	if !this.LocallyManageTransactions {
		return
	}
	err := this.OrmInstance.Rollback()
	if err != nil {
		panic(err)
	}
}

func (this *OrmWrapper) RollbackTransaction_NoPanic() {
	defer func() {
		recover()
	}()
	this.OrmInstance.Rollback()
}

func (this *OrmWrapper) CommitTransaction() {
	if !this.LocallyManageTransactions {
		return
	}
	err := this.OrmInstance.Commit()
	if err != nil {
		panic(err)
	}
}
