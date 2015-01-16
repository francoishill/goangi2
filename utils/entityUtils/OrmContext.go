package entityUtils

import (
	. "github.com/francoishill/goangi2/utils/debugUtils"
	. "github.com/francoishill/goangi2/utils/loggingUtils"
)

type OrmContext struct {
	Logger     ILogger
	OrmWrapper *OrmWrapper
}

func CreateOrmContext(logger ILogger, possibleParentTransactionalOrmWrapper *OrmWrapper, beginTransaction bool) *OrmContext {
	loggerToUse := logger
	if loggerToUse == nil {
		loggerToUse = CreateNewFmtLogger()
	}
	if possibleParentTransactionalOrmWrapper == nil {
		possibleParentTransactionalOrmWrapper = CreateNewOrmWrapper(nil)
	}
	ormCtx := &OrmContext{
		Logger:     loggerToUse,
		OrmWrapper: possibleParentTransactionalOrmWrapper,
	}
	if beginTransaction {
		ormCtx.OrmWrapper.BeginTransaction()
	}
	return ormCtx
}

func CreateDefaultOrmContext() *OrmContext {
	var ormWrapper *OrmWrapper = nil
	beginTransaction := false
	return CreateOrmContext(CreateNewFmtLogger(), ormWrapper, beginTransaction)
}

func (this *OrmContext) getCurrentStackTrace() string {
	return GetFullStackTrace_Pretty()
}

func (this *OrmContext) CommitOnSuccessOrRollbackOnPanic() {
	if r := recover(); r != nil {
		this.OrmWrapper.RollbackTransaction_NoPanic()
		panic(r)
	} else {
		this.OrmWrapper.CommitTransaction()
	}
}

func (this *OrmContext) RecoverAndPrintIfPanic_Error(additionalMessageNoFullStop string, funcOnErrorCatched func(interface{}), funcOnFinally func()) {
	if r := recover(); r != nil {
		this.Logger.Error("%s. Error: %+v. Stack: %s", additionalMessageNoFullStop, r, this.getCurrentStackTrace())
		if funcOnErrorCatched != nil {
			funcOnErrorCatched(r)
		}
	}
	if funcOnFinally != nil {
		funcOnFinally()
	}
}
