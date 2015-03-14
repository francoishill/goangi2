package entityUtils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type beegoOrmRepo struct {
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func (this *beegoOrmRepo) CheckEntityExistsWithPK(ormContext *OrmContext, entityObj interface{}) bool {
	ormContextToUse := CreateOrmContext_FromAnother(ormContext, false)
	err := ormContextToUse.OrmWrapper.OrmInstance.Read(entityObj)
	if this.BaseErrorIsZeroRowsFound(err) {
		return false
	}
	checkError(err)
	return true
}

func (this *beegoOrmRepo) BaseReadEntityUsingPK(ormContext *OrmContext, entityObj interface{}, relatedFieldsToLoad *RelatedFieldsToLoad) {
	ormContextToUse := CreateOrmContext_FromAnother(ormContext, false)
	err := ormContextToUse.OrmWrapper.OrmInstance.Read(entityObj)
	checkError(err)

	if relatedFieldsToLoad != nil {
		allRelatedFields := relatedFieldsToLoad.GetFieldNames(true, true)
		for _, fieldName := range allRelatedFields {
			OrmRepo.BaseLoadRelatedFields(ormContext, entityObj, fieldName)
		}
	}
}

func (this *beegoOrmRepo) BaseReadEntityUsingFields(ormContext *OrmContext, entityObj interface{}, relatedFieldsToLoad *RelatedFieldsToLoad, fields ...string) (success bool) {
	ormContextToUse := CreateOrmContext_FromAnother(ormContext, false)
	err := ormContextToUse.OrmWrapper.OrmInstance.Read(entityObj, fields...)
	if this.BaseErrorIsZeroRowsFound(err) {
		return false
	}
	checkError(err)

	if relatedFieldsToLoad != nil {
		allRelatedFields := relatedFieldsToLoad.GetFieldNames(true, true)
		for _, fieldName := range allRelatedFields {
			OrmRepo.BaseLoadRelatedFields(ormContext, entityObj, fieldName)
		}
	}

	return true
}

func (this *beegoOrmRepo) CheckIfExistByField(ormContext *OrmContext, tableName string, field string, value interface{}, skipId int64, caseSensitive bool) bool {
	qs := orm.NewOrm().
		QueryTable(tableName).
		OrderBy("-Id")

	if caseSensitive {
		qs = qs.Filter(field, value)
	} else {
		qs = qs.Filter(field+"__iexact", value)
	}
	if skipId > 0 {
		qs = qs.Exclude("Id", skipId)
	}
	return qs.Exist()
}

func (this *beegoOrmRepo) BaseInsertEntity(ormContext *OrmContext, entityObj interface{}) {
	var ormWrapperToUse *OrmWrapper
	if ormContext == nil {
		ormContext = CreateDefaultOrmContext()
	}
	ormWrapperToUse = CreateNewOrmWrapper(ormContext.OrmWrapper)
	hasCalledBeginTransaction := false
	defer ormContext.RecoverAndPrintIfPanic_Error(fmt.Sprintf("Unable to insert entity (%+v)", entityObj), func(recoveryObj interface{}) {
		if hasCalledBeginTransaction {
			ormWrapperToUse.RollbackTransaction()
		}
		panic(recoveryObj)
	}, nil)

	ormWrapperToUse.BeginTransaction()
	hasCalledBeginTransaction = true

	_, err := ormWrapperToUse.OrmInstance.Insert(entityObj)
	checkError(err)

	ormWrapperToUse.CommitTransaction()
}

func (this *beegoOrmRepo) BaseUpdateEntity(ormContext *OrmContext, baseEntity interface{}, modifiedEntity interface{}, onlyAllowTheseFieldsToSave ...string) []ChangedField {
	changedFields := FormChanges(baseEntity, modifiedEntity, nil, onlyAllowTheseFieldsToSave)
	if len(changedFields) == 0 {
		return []ChangedField{}
	}
	SetFormValues(modifiedEntity, baseEntity, nil, onlyAllowTheseFieldsToSave)

	var filteredChangedFields []ChangedField
	if len(onlyAllowTheseFieldsToSave) == 0 {
		filteredChangedFields = changedFields
	} else {
		for _, fld := range changedFields {
			fieldName := fld.FieldName
			for _, fn := range onlyAllowTheseFieldsToSave {
				if strings.EqualFold(fieldName, fn) {
					filteredChangedFields = append(filteredChangedFields, fld)
					break
				}
			}
		}
	}

	changedFieldNames := []string{}
	if DoesStructContainField(baseEntity, "Updated") {
		changedFieldNames = append(changedFieldNames, "Updated")
	}
	for _, fld := range filteredChangedFields {
		changedFieldNames = append(changedFieldNames, fld.FieldName)
	}

	var ormWrapperToUse *OrmWrapper
	if ormContext == nil {
		ormContext = CreateDefaultOrmContext()
	}
	ormWrapperToUse = CreateNewOrmWrapper(ormContext.OrmWrapper)

	hasCalledBeginTransaction := false
	defer ormContext.RecoverAndPrintIfPanic_Error(fmt.Sprintf("Unable to update entity (orig: %+v, modified: %+v)", baseEntity, modifiedEntity), func(recoveryObj interface{}) {
		if hasCalledBeginTransaction {
			ormWrapperToUse.RollbackTransaction()
		}
		panic(recoveryObj)
	}, nil)

	ormWrapperToUse.BeginTransaction()
	hasCalledBeginTransaction = true

	_, err := ormWrapperToUse.OrmInstance.Update(baseEntity, changedFieldNames...)
	checkError(err)

	ormWrapperToUse.CommitTransaction()
	return filteredChangedFields
}

func (this *beegoOrmRepo) BaseDeleteEntity(ormContext *OrmContext, entityObj interface{}) {
	var ormWrapperToUse *OrmWrapper
	if ormContext == nil {
		ormContext = CreateDefaultOrmContext()
	}
	ormWrapperToUse = CreateNewOrmWrapper(ormContext.OrmWrapper)

	hasCalledBeginTransaction := false
	defer ormContext.RecoverAndPrintIfPanic_Error(fmt.Sprintf("Unable to delete entity (%+v)", entityObj), func(recoveryObj interface{}) {
		if hasCalledBeginTransaction {
			ormWrapperToUse.RollbackTransaction()
		}
		panic(recoveryObj)
	}, nil)

	ormWrapperToUse.BeginTransaction()
	hasCalledBeginTransaction = true

	_, err := ormWrapperToUse.OrmInstance.Delete(entityObj)
	checkError(err)

	ormWrapperToUse.CommitTransaction()
}

func (this *beegoOrmRepo) BaseUpdateM2MByAddAndRemove(ormContext *OrmContext,
	entityObj interface{}, columnNameOfRelationship string,
	removeListToRelationship, addListToRelationship []interface{}) {

	var err error

	var ormWrapperToUse *OrmWrapper
	if ormContext == nil {
		ormContext = CreateDefaultOrmContext()
	}
	ormWrapperToUse = CreateNewOrmWrapper(ormContext.OrmWrapper)

	hasCalledBeginTransaction := false
	defer ormContext.RecoverAndPrintIfPanic_Error(
		fmt.Sprintf("Unable to UpdateM2MByAddAndRemove (entity: %+v, colname: %s, removelist: %+v, addlist: %+v)", entityObj, columnNameOfRelationship, removeListToRelationship, addListToRelationship),
		func(recoveryObj interface{}) {
			if hasCalledBeginTransaction {
				ormWrapperToUse.RollbackTransaction()
			}
			panic(recoveryObj)
		}, nil)

	m2mObj := ormWrapperToUse.OrmInstance.QueryM2M(entityObj, columnNameOfRelationship)

	ormWrapperToUse.BeginTransaction()
	hasCalledBeginTransaction = true

	if len(removeListToRelationship) > 0 {
		_, err = m2mObj.Remove(removeListToRelationship...)
		checkError(err)
	}

	if len(addListToRelationship) > 0 {
		_, err = m2mObj.Add(addListToRelationship...)
		checkError(err)
	}

	ormWrapperToUse.CommitTransaction()
}

// func (this *beegoOrmRepo) BaseListEntities_OrderBy_Limit_Offset__WhereM2MCountIsZero(
// 	ormContext *OrmContext,
// 	dbDriverName string,
// 	leftTableName string,
// 	leftTable_IdField string,
// 	m2mTableName string,
// 	m2mTableName_LeftTableIdField string,
// 	columnNameOfRelationship string,
// 	otherTableName string,
// 	orderByFields []string,
// 	limit int64,
// 	offset int64,
// 	relatedFieldsToLoad *RelatedFieldsToLoad,
// 	sliceToPopulatePointer interface{}) {

// 	/*qs := getNewBeegoOrm().QueryTable(leftTableName)
// 	if limit > 0 && offset > 0 {
// 		qs = qs.Limit(limit, offset)
// 	} else if limit > 0 {
// 		qs = qs.Limit(limit)
// 	} else if offset > 0 {
// 		qs = qs.Offset(offset)
// 	}*/

// 	qb, err := orm.NewQueryBuilder(dbDriverName)
// 	if err != nil {
// 		panic(err)
// 	}
// 	qb.Select(leftTableName + "*").
// 		From(leftTableName).
// 		LeftJoin(m2mTableName).On(leftTableName + "." + leftTable_IdField + " = " + m2mTableName + "." + m2mTableName_LeftTableIdField)

// 	if len(orderByFields) > 0 {
// 		qb = qb.OrderBy(orderByFields...)
// 	}

// 	if limit > 0 {
// 		qb = qb.Limit(limit)
// 	}

// 	if offset > 0 {
// 		qb = qb.Offset(offset)
// 	}

// 	sql := qb.String()

// 	if relatedFieldsToLoad != nil {
// 		fieldNames := relatedFieldsToLoad.GetFieldNames(true, false)
// 		for _, relFieldName := range fieldNames {
// 			qs = qs.RelatedSel(relFieldName)
// 		}
// 	}

// 	o := orm.NewOrm()
// 	o.Raw(sql, 20).QueryRows(&users)
// }

func (this *beegoOrmRepo) BaseCountM2M(
	ormContext *OrmContext,
	entityObj interface{},
	columnNameOfRelationship string) int64 {
	var ormWrapperToUse *OrmWrapper
	if ormContext == nil {
		ormContext = CreateDefaultOrmContext()
	}
	ormWrapperToUse = CreateNewOrmWrapper(ormContext.OrmWrapper)

	m2mObj := ormWrapperToUse.OrmInstance.QueryM2M(entityObj, columnNameOfRelationship)
	cnt, err := m2mObj.Count()
	checkError(err)
	return cnt
}

func (this *beegoOrmRepo) BaseM2MRelationExists(
	ormContext *OrmContext,
	entityObj interface{},
	columnNameOfRelationship string,
	relationEntity interface{}) bool {
	var ormWrapperToUse *OrmWrapper
	if ormContext == nil {
		ormContext = CreateDefaultOrmContext()
	}
	ormWrapperToUse = CreateNewOrmWrapper(ormContext.OrmWrapper)

	m2mObj := ormWrapperToUse.OrmInstance.QueryM2M(entityObj, columnNameOfRelationship)
	return m2mObj.Exist(relationEntity)
}

func (this *beegoOrmRepo) BaseListEntities_ANDFilters_OrderBy(
	ormContext *OrmContext,
	queryTableName string,
	fieldFilters []map[string]interface{},
	orderByFields []string,
	relatedFieldsToLoad *RelatedFieldsToLoad,
	sliceToPopulatePointer interface{}) {

	limit := int64(DEFAULT_QUERY_LIMIT)
	offset := int64(0) //default offset always 0
	this.BaseListEntities_ANDFilters_OrderBy_Limit_Offset(
		ormContext, queryTableName, fieldFilters, orderByFields, limit, offset, relatedFieldsToLoad, sliceToPopulatePointer)
}

func (this *beegoOrmRepo) BaseListEntities_ANDFilters_OrderBy_Limit_Offset(
	ormContext *OrmContext,
	queryTableName string,
	fieldFilters []map[string]interface{},
	orderByFields []string,
	limit int64,
	offset int64,
	relatedFieldsToLoad *RelatedFieldsToLoad,
	sliceToPopulatePointer interface{}) {

	ormContextToUse := CreateOrmContext_FromAnother(ormContext, false)
	qs := ormContextToUse.OrmWrapper.OrmInstance.QueryTable(queryTableName)
	if limit > 0 && offset > 0 {
		qs = qs.Limit(limit, offset)
	} else if limit > 0 {
		qs = qs.Limit(limit)
	} else if offset > 0 {
		qs = qs.Offset(offset)
	}

	if fieldFilters != nil && len(fieldFilters) > 0 {
		//Each element of the array of fieldFilters can contain a list of OR filters on a field
		cond := orm.NewCondition()
		andConditionList := []*orm.Condition{}
		for ind, _ := range fieldFilters {
			tmpCond := cond
			for fieldName, fieldVal := range fieldFilters[ind] {
				tmpCond = tmpCond.Or(fieldName, fieldVal)
			}
			andConditionList = append(andConditionList, tmpCond)
		}

		finalCondition := cond
		for ind, _ := range andConditionList {
			finalCondition = finalCondition.AndCond(andConditionList[ind])
		}

		qs = qs.SetCond(finalCondition)
	}

	if relatedFieldsToLoad != nil {
		fieldNames := relatedFieldsToLoad.GetFieldNames(true, false)
		for _, relFieldName := range fieldNames {
			qs = qs.RelatedSel(relFieldName)
		}
	}

	if len(orderByFields) > 0 {
		qs = qs.OrderBy(orderByFields...)
	}

	var sliceValToUse reflect.Value

	switch reflect.TypeOf(sliceToPopulatePointer).Kind() {
	case reflect.Ptr:
		sliceVal := reflect.Indirect(reflect.ValueOf(sliceToPopulatePointer))

		switch reflect.TypeOf(sliceVal.Interface()).Kind() {
		case reflect.Slice:
			sliceValToUse = sliceVal
		}
		break
	case reflect.Slice:
		sliceVal := reflect.ValueOf(sliceToPopulatePointer)
		sliceValToUse = sliceVal
		break
	default:
		panic("Unexpected slice type to list entities")
		break
	}

	_, err := dupListObjects(qs, sliceValToUse.Addr().Interface())
	checkError(err)

	if relatedFieldsToLoad != nil {
		externalFields := relatedFieldsToLoad.GetFieldNames(false, true)

		for _, relFieldName := range externalFields {
			for i := 0; i < sliceValToUse.Len(); i++ {
				//sliceEntityPointer := sliceValToUse.Index(i).Addr().Interface()
				sliceEntityPointer := sliceValToUse.Index(i).Interface()
				OrmRepo.BaseLoadRelatedFields(ormContext, sliceEntityPointer, relFieldName)
			}
		}
	}
}

func (this *beegoOrmRepo) BaseCountEntities_ANDFilters(ormContext *OrmContext, queryTableName string, fieldFilters []map[string]interface{}) int64 {
	ormContextToUse := CreateOrmContext_FromAnother(ormContext, false)
	qs := ormContextToUse.OrmWrapper.OrmInstance.QueryTable(queryTableName)

	if fieldFilters != nil && len(fieldFilters) > 0 {
		//Each element of the array of fieldFilters can contain a list of OR filters on a field
		cond := orm.NewCondition()
		andConditionList := []*orm.Condition{}
		for ind, _ := range fieldFilters {
			tmpCond := cond
			for fieldName, fieldVal := range fieldFilters[ind] {
				tmpCond = tmpCond.Or(fieldName, fieldVal)
			}
			andConditionList = append(andConditionList, tmpCond)
		}

		finalCondition := cond
		for ind, _ := range andConditionList {
			finalCondition = finalCondition.AndCond(andConditionList[ind])
		}

		qs = qs.SetCond(finalCondition)
	}

	cnt, err := qs.Count()
	checkError(err)

	return cnt
}

func (this *beegoOrmRepo) BaseLoadRelatedFields(ormContext *OrmContext, m interface{}, fieldName string) int64 {
	var numLoaded int64
	var err error

	ormContextToUse := CreateOrmContext_FromAnother(ormContext, false)

	numLoaded, err = ormContextToUse.OrmWrapper.OrmInstance.LoadRelated(m, fieldName)
	if err != nil {
		//If the error is No Rows Found it means probably that we are trying to get a FK relationship where the field is NULL
		if err != orm.ErrNoRows {
			checkError(err)
		}
	}
	return numLoaded
}

func (this *beegoOrmRepo) BaseErrorIsZeroRowsFound(err error) bool {
	return err == orm.ErrNoRows
}

func getNewBeegoOrm() orm.Ormer {
	return orm.NewOrm()
}

func dupGetIdOfItem(baseEntity interface{}) int64 {
	valOfItem := reflect.ValueOf(baseEntity).Elem()
	typeOfItem := valOfItem.Type()
	var itemId int64
	if _, gotField := typeOfItem.FieldByName("Id"); !gotField {
		panic("Cannot find Id field to use for IRevisionableEntity")
	} else if idInt64, gotInt64Val := valOfItem.FieldByName("Id").Interface().(int64); !gotInt64Val {
		panic("Cannot convert Id field to int64 for IRevisionableEntity")
	} else {
		itemId = idInt64
	}
	return itemId
}

func dupListObjects(qs orm.QuerySeter, objs interface{}, onlyReturnTheseFields ...string) (int64, error) {
	nums, err := qs.All(objs, onlyReturnTheseFields...)
	if err != nil {
		return 0, err
	}
	return nums, err
}
