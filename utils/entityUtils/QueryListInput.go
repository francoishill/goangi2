package entityUtils

type QueryListInput struct {
	Limit         int64
	Offset        int64
	QueryFilter   *QueryFilter
	RelatedFields *RelatedFieldsToLoad
	OrderByFields []string
}

func QueryListInput_Construct(limit, offset int64, queryFilter *QueryFilter, relatedFields *RelatedFieldsToLoad, orderByFields ...string) *QueryListInput {
	return &QueryListInput{
		Limit:         limit,
		Offset:        offset,
		QueryFilter:   queryFilter,
		RelatedFields: relatedFields,
		OrderByFields: orderByFields,
	}
}

func (this *QueryListInput) ExecuteQuery(ormContext *OrmContext, tableName string, sliceToPopulatePointer interface{}) {
	if this.Limit > 0 || this.Offset > 0 {
		OrmRepo.BaseListEntities_ANDFilters_OrderBy_Limit_Offset(
			ormContext, tableName, this.QueryFilter, this.OrderByFields, this.Limit, this.Offset, this.RelatedFields, sliceToPopulatePointer)
	} else {
		OrmRepo.BaseListEntities_ANDFilters_OrderBy(
			ormContext, tableName, this.QueryFilter, this.OrderByFields, this.RelatedFields, sliceToPopulatePointer)
	}
}
