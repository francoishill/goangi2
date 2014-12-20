package entityUtils

func CreateRelatedField(fieldName string, canRelatedSel bool) *RelatedField {
	return &RelatedField{
		FieldName:     fieldName,
		CanRelatedSel: canRelatedSel,
	}
}

func CreateRelatedFieldsToLoad(relatedFields ...*RelatedField) *RelatedFieldsToLoad {
	return &RelatedFieldsToLoad{
		fields: relatedFields,
	}
}

type RelatedField struct {
	FieldName string

	//CanRelatedSel is TRUE for
	// ForeignKey
	// Rel(one)
	//But FALSE for
	// Rev(one)
	// Rev Many
	// M2M
	CanRelatedSel bool
}

type RelatedFieldsToLoad struct {
	fields []*RelatedField
}

func (this *RelatedFieldsToLoad) AppendIfTrue(checkIfTrue bool, relatedFieldToAppend *RelatedField) {
	if !checkIfTrue {
		return
	}
	this.fields = append(this.fields, relatedFieldToAppend)
}

func (this *RelatedFieldsToLoad) GetFieldNames(includeCanRelatedSel, includeNOTcanRelatedSel bool) []string {
	fieldNames := []string{}
	for _, rl := range this.fields {
		if !includeCanRelatedSel && rl.CanRelatedSel {
			continue
		}
		if !includeNOTcanRelatedSel && !rl.CanRelatedSel {
			continue
		}
		fieldNames = append(fieldNames, rl.FieldName)
	}
	return fieldNames
}
