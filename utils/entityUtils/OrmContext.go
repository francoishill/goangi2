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
	ormCtx := &OrmContext{
		Logger:     loggerToUse,
		OrmWrapper: CreateNewOrmWrapper(possibleParentTransactionalOrmWrapper),
	}
	if beginTransaction {
		ormCtx.OrmWrapper.BeginTransaction()
	}
	return ormCtx
}

func CreateOrmContext_FromAnother(ormContext *OrmContext, beginTransaction bool) *OrmContext {
	if ormContext == nil {
		returnContext := CreateDefaultOrmContext()
		if beginTransaction {
			returnContext.OrmWrapper.BeginTransaction()
		}
		return returnContext
	} else {
		return CreateOrmContext(ormContext.Logger, ormContext.OrmWrapper, beginTransaction)
	}
}

func CreateOrmContext_FromAnother_ButCreateNewOrmWrapper(ormContext *OrmContext, beginTransaction bool) *OrmContext {
	if ormContext == nil {
		returnContext := CreateDefaultOrmContext()
		if beginTransaction {
			returnContext.OrmWrapper.BeginTransaction()
		}
		return returnContext
	} else {
		return CreateOrmContext(ormContext.Logger, CreateNewOrmWrapper(nil), beginTransaction)
	}
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
