package entityUtils

type iOrmRepo interface {
	BaseReadEntityUsingPK(ormContext *OrmContext, entityObj interface{}, relatedFieldsToLoad *RelatedFieldsToLoad)

	BaseReadEntityUsingFields(ormContext *OrmContext, entityObj interface{}, relatedFieldsToLoad *RelatedFieldsToLoad, fields ...string) bool

	CheckIfExistByField(ormContext *OrmContext, tableName, field string, value interface{}, skipId int64, caseSensitive bool) bool

	BaseInsertEntity(
		ormContext *OrmContext,
		entityObj interface{})

	BaseUpdateEntity(
		ormContext *OrmContext,
		baseEntity interface{},
		modifiedEntity interface{},
		onlyAllowTheseFieldsToSave ...string) []ChangedField

	BaseDeleteEntity(
		ormContext *OrmContext,
		entityObj interface{})

	BaseUpdateM2MByAddAndRemove(
		ormContext *OrmContext,
		entityObj interface{},
		columnNameOfRelationship string,
		removeListToRelationship,
		addListToRelationship []interface{})

	/*BaseListEntities_OrderBy_Limit_Offset__WhereM2MCountIsZero(
	ormContext *OrmContext,
	queryTableName string,
	columnNameOfRelationship string,
	orderByFields []string,
	limit int64,
	offset int64,
	relatedFieldsToLoad *RelatedFieldsToLoad,
	sliceToPopulatePointer interface{})*/

	BaseCountM2M(
		ormContext *OrmContext,
		entityObj interface{},
		columnNameOfRelationship string) int64

	BaseM2MRelationExists(
		ormContext *OrmContext,
		entityObj interface{},
		columnNameOfRelationship string,
		relationEntity interface{}) bool

	BaseListEntities_ANDFilters_OrderBy(
		ormContext *OrmContext,
		queryTableName string,
		fieldFilters []map[string]interface{}, //Map of FieldName + FieldValue
		orderByFields []string,
		relatedFieldsToLoad *RelatedFieldsToLoad,
		sliceToPopulatePointer interface{})

	BaseListEntities_ANDFilters_OrderBy_Limit_Offset(
		ormContext *OrmContext,
		queryTableName string,
		fieldFilters []map[string]interface{}, //Map of FieldName + FieldValue
		orderByFields []string,
		limit int64,
		offset int64,
		relatedFieldsToLoad *RelatedFieldsToLoad,
		sliceToPopulatePointer interface{})

	BaseCountEntities_ANDFilters(
		ormContext *OrmContext,
		queryTableName string,
		fieldFilters []map[string]interface{}) int64

	BaseLoadRelatedFields(ormContext *OrmContext, m interface{}, fieldName string) int64

	BaseErrorIsZeroRowsFound(err error) bool
}

var (
	OrmRepo = iOrmRepo(new(beegoOrmRepo))
)
