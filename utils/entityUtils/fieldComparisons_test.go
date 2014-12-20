package entityUtils

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

type tmpPersonStructForSetFormValues struct {
	Name     string
	Age      int
	Employed bool
}

func doTestPassingNonStructPointer_ForTestSetFormValues(t *testing.T, obj1, obj2 interface{}, extraMessage string) {
	defer func() {
		recover()
	}()

	SetFormValues(obj1, obj2, nil, nil)

	AssertAndFailNow(t, false, "The previous line should not have succeeded. We should require both to be POINTERS, but we are not passing in a pointer. "+extraMessage)
}
func TestSetFormValues(t *testing.T) {
	tmpStruct1 := tmpPersonStructForSetFormValues{}
	tmpStruct2 := tmpPersonStructForSetFormValues{}
	doTestPassingNonStructPointer_ForTestSetFormValues(t, tmpStruct1, tmpStruct2, "Case 1")  //Both non pointers
	doTestPassingNonStructPointer_ForTestSetFormValues(t, &tmpStruct1, tmpStruct2, "Case 2") //Only one not pointer
	doTestPassingNonStructPointer_ForTestSetFormValues(t, tmpStruct1, &tmpStruct2, "Case 3") //Only one not pointer

	var destinationPersonToUpdate *tmpPersonStructForSetFormValues

	Convey("Having a source and destination person", t, func() {
		These tests will not pass, have to continue changing all `AssertEqualsStringAnd....` related code to goconvey code

		sourcePersonToCopyFrom := &tmpPersonStructForSetFormValues{
			Name:     "NewName",
			Age:      26,
			Employed: false,
		}
		destinationPersonToUpdate = &tmpPersonStructForSetFormValues{
			Name:     "OrigName",
			Age:      25,
			Employed: true,
		}
		AssertEqualStringAndFailNow(t, destinationPersonToUpdate.Name, "OrigName", "Unexpected Name (case 1) for destinationPersonToUpdate.")
		AssertEqualIntAndFailNow(t, destinationPersonToUpdate.Age, 25, "Unexpected Age (case 1) for destinationPersonToUpdate.")
		SetFormValues(sourcePersonToCopyFrom, destinationPersonToUpdate, nil, nil) //Transfer from sourcePersonToCopyFrom to destinationPersonToUpdate
		AssertEqualStringAndFailNow(t, destinationPersonToUpdate.Name, "NewName", "Name (case 1) did not change with SetFormValues.")
		AssertEqualIntAndFailNow(t, destinationPersonToUpdate.Age, 26, "Age (case 1) did not change with SetFormValues.")
	})

	//Reset destinationPersonToUpdate for this case
	sourceAnonToCopyFrom := &struct {
		Name string
		Age  int
	}{
		Name: "NewNameAnon",
		Age:  27,
	}
	destinationPersonToUpdate = &tmpPersonStructForSetFormValues{
		Name:     "OrigName",
		Age:      25,
		Employed: true,
	}
	AssertEqualStringAndFailNow(t, destinationPersonToUpdate.Name, "OrigName", "Unexpected Name (case 2) for destinationPersonToUpdate.")
	AssertEqualIntAndFailNow(t, destinationPersonToUpdate.Age, 25, "Unexpected Age (case 2) for destinationPersonToUpdate.")
	SetFormValues(sourceAnonToCopyFrom, destinationPersonToUpdate, nil, nil)
	AssertEqualStringAndFailNow(t, destinationPersonToUpdate.Name, "NewNameAnon", "Name (case 2) did not change with SetFormValues.")
	AssertEqualIntAndFailNow(t, destinationPersonToUpdate.Age, 27, "Age (case 2) did not change with SetFormValues.")

	//Reset destinationPersonToUpdate for this case
	sourceAnonToCopyFrom2 := &struct {
		Name string
		Age  int
	}{
		Name: "NewNameAnon",
		Age:  27,
	}
	destinationPersonToUpdate = &tmpPersonStructForSetFormValues{
		Name:     "OrigName",
		Age:      25,
		Employed: true,
	}
	AssertEqualStringAndFailNow(t, destinationPersonToUpdate.Name, "OrigName", "Unexpected Name (case 3) for destinationPersonToUpdate.")
	AssertEqualIntAndFailNow(t, destinationPersonToUpdate.Age, 25, "Unexpected Age (case 3) for destinationPersonToUpdate.")

	SetFormValues(sourceAnonToCopyFrom2, destinationPersonToUpdate, []string{"Name", "Age"}, nil) //Skip both Name+Age
	AssertEqualStringAndFailNow(t, destinationPersonToUpdate.Name, "OrigName", "Name (case 3) was NOT supposed to update with SetFormValues as we ignored Name+Age.")
	AssertEqualIntAndFailNow(t, destinationPersonToUpdate.Age, 25, "Age (case 3) was NOT supposed to update with SetFormValues as we ignored Name+Age.")

	SetFormValues(sourceAnonToCopyFrom2, destinationPersonToUpdate, []string{"Name"}, nil) //Skip only Name
	AssertEqualStringAndFailNow(t, destinationPersonToUpdate.Name, "OrigName", "Name (case 4) was NOT supposed to update with SetFormValues as we ignored Name.")
	AssertEqualIntAndFailNow(t, destinationPersonToUpdate.Age, 27, "Age (case 4) did not change with SetFormValues.")

	//Just reset the values of Name+Age
	destinationPersonToUpdate = &tmpPersonStructForSetFormValues{
		Name:     "OrigName",
		Age:      25,
		Employed: true,
	}
	SetFormValues(sourceAnonToCopyFrom2, destinationPersonToUpdate, []string{"Age"}, nil) //Skip only Age
	AssertEqualStringAndFailNow(t, destinationPersonToUpdate.Name, "NewNameAnon", "Name (case 4) did not change with SetFormValues.")
	AssertEqualIntAndFailNow(t, destinationPersonToUpdate.Age, 25, "Age (case 4) was NOT supposed to update with SetFormValues as we ignored Age.")
}

type tmpPersonStructForFormChanges struct {
	Name           string
	SurnameIgnored string `form:"-"` //This field should be ignored in field comparisons
	Age            int
	Employed       bool

	Created time.Time `orm:"auto_now_add"`
	Updated time.Time `orm:"auto_now"`
}

func TestFormChanges(t *testing.T) {
	now := time.Now()

	originalPerson := &tmpPersonStructForFormChanges{
		Name:           "Name",
		SurnameIgnored: "Surname",
		Age:            24,
		Employed:       true,
		Created:        now,
		Updated:        now,
	}

	//Modify nothing
	modifiedPerson1 := &tmpPersonStructForFormChanges{
		Name:           "Name",
		SurnameIgnored: "Surname",
		Age:            24,
		Employed:       true,
		Created:        now,
		Updated:        now,
	}
	changedFields1 := FormChanges(originalPerson, modifiedPerson1, nil, nil)
	AssertEqualIntAndFailNow(t, len(changedFields1), 0, "Unexpected number of changed fields for identical entity.")

	//Only modify Name
	modifiedPerson2 := &tmpPersonStructForFormChanges{
		Name:           "NewName",
		SurnameIgnored: "Surname",
		Age:            24,
		Employed:       true,
		Created:        now,
		Updated:        now,
	}
	changedFields2 := FormChanges(originalPerson, modifiedPerson2, nil, nil)
	AssertEqualIntAndFailNow(t, len(changedFields2), 1, "Unexpected number of changed fields for only modifying only the Name field.")
	changedField2_0 := changedFields2[0]
	AssertEqualStringAndFailNow(t, changedField2_0.FieldName, "Name", "Unexpected changed FieldName.")
	AssertEqualStringAndFailNow(t, changedField2_0.OldValue, "Name", "Unexpected changed OldValue.")
	AssertEqualStringAndFailNow(t, changedField2_0.NewValue, "NewName", "Unexpected changed NewValue.")

	//Only modify Age
	modifiedPerson3 := &tmpPersonStructForFormChanges{
		Name:           "Name",
		SurnameIgnored: "Surname",
		Age:            25,
		Employed:       true,
		Created:        now,
		Updated:        now,
	}
	changedFields3 := FormChanges(originalPerson, modifiedPerson3, nil, nil)
	AssertEqualIntAndFailNow(t, len(changedFields3), 1, "Unexpected number of changed fields for only modifying only the Age field.")
	changedField3_0 := changedFields3[0]
	AssertEqualStringAndFailNow(t, changedField3_0.FieldName, "Age", "Unexpected changed FieldName.")
	AssertEqualStringAndFailNow(t, changedField3_0.OldValue, "24", "Unexpected changed OldValue.")
	AssertEqualStringAndFailNow(t, changedField3_0.NewValue, "25", "Unexpected changed NewValue.")

	//Only modify Employed
	modifiedPerson4 := &tmpPersonStructForFormChanges{
		Name:           "Name",
		SurnameIgnored: "Surname",
		Age:            24,
		Employed:       false,
		Created:        now,
		Updated:        now,
	}
	changedFields4 := FormChanges(originalPerson, modifiedPerson4, nil, nil)
	AssertEqualIntAndFailNow(t, len(changedFields4), 1, "Unexpected number of changed fields for only modifying only the Employed field.")
	changedField4_0 := changedFields4[0]
	AssertEqualStringAndFailNow(t, changedField4_0.FieldName, "Employed", "Unexpected changed FieldName.")
	AssertEqualStringAndFailNow(t, changedField4_0.OldValue, "true", "Unexpected changed OldValue.")
	AssertEqualStringAndFailNow(t, changedField4_0.NewValue, "false", "Unexpected changed NewValue.")

	//Modify all fields
	modifiedPerson5 := &tmpPersonStructForFormChanges{
		Name:           "NewName",
		SurnameIgnored: "NewSurname",
		Age:            26,
		Employed:       false,
		Created:        now,
		Updated:        now,
	}
	changedFields5a := FormChanges(originalPerson, modifiedPerson5, nil, nil)
	AssertEqualIntAndFailNow(t, len(changedFields5a), 3, "Unexpected number of changed fields for modifying Name, Age and Employed (Surname ignored) fields.")
	changedFields5b := FormChanges(originalPerson, modifiedPerson5, []string{"Name", "Age"}, nil) //Skipping the changes of Name+Age changes
	AssertEqualIntAndFailNow(t, len(changedFields5b), 1, "Unexpected number of changed fields for modifying Name, Age and Employed (Surname ignored and skipping Name+Age) fields.")

	//Modify nothing but only the Created and Updated, but these two fields should be ignored.
	modifiedPerson6 := &tmpPersonStructForFormChanges{
		Name:           "Name",
		SurnameIgnored: "Surname",
		Age:            24,
		Employed:       true,
		Created:        now.Add(10 * time.Hour),
		Updated:        now.Add(10 * time.Hour),
	}
	changedFields6 := FormChanges(originalPerson, modifiedPerson6, nil, nil)
	AssertEqualIntAndFailNow(t, len(changedFields6), 0, "Unexpected number of changed fields for entity where ONLY the Created and Updated dates changed. These dates are expected to be ignored.")

	//Modify nothing but only the form ignored field
	modifiedPerson7 := &tmpPersonStructForFormChanges{
		Name:           "Name",
		SurnameIgnored: "NewSurname",
		Age:            24,
		Employed:       true,
		Created:        now,
		Updated:        now,
	}
	changedFields7 := FormChanges(originalPerson, modifiedPerson7, nil, nil)
	AssertEqualIntAndFailNow(t, len(changedFields7), 0, "Unexpected number of changed fields for entity where ONLY the SurnameIgnored changed, which is a field to be ignored..")
}

type tmpEmptyStruct1 struct{}
type tmpEmptyStruct2 struct{}

func doTestPassingNonStructPointer_ForTestGetAllFieldsOfStruct(t *testing.T) {
	defer func() {
		recover()
	}()

	tmpStruct1 := tmpEmptyStruct1{}
	ignoreSettings := &IgnoreFieldTypes{OrmSkipped: false, Created: false, Updated: false, Relationals: false}
	GetAllFieldsOfStruct(tmpStruct1, ignoreSettings)

	AssertAndFailNow(t, false, "The previous line should not have succeeded. We should require a POINTER, but we are not passing in a pointer.")
}
func TestGetAllFieldsOfStruct_AndGetAllFieldNamesOfStruct(t *testing.T) {
	doTestPassingNonStructPointer_ForTestGetAllFieldsOfStruct(t)

	person2 := struct {
		Name               string
		Age                int
		Employed           bool
		PointerField       *tmpEmptyStruct1 // Pointer fields are ignored by default
		RelField2          *tmpEmptyStruct2 `orm:"rel(fk)"` // Pointer fields are ignored by default
		TempCodeOrmSkipped string           `orm:"-"`
		Created            time.Time        `orm:"auto_now_add"`
		Updated            time.Time        `orm:"auto_now"`
	}{
		Name:               "Name",
		Age:                25,
		Employed:           true,
		TempCodeOrmSkipped: "123lkjljks",
		Created:            time.Now(),
		Updated:            time.Now(),
	}
	ignoreSettingsNone := &IgnoreFieldTypes{OrmSkipped: false, Created: false, Updated: false, Relationals: false}
	allFields := GetAllFieldsOfStruct(&person2, ignoreSettingsNone)
	allFieldNames := GetAllFieldNamesOfStruct(&person2, ignoreSettingsNone)
	AssertEqualIntAndFailNow(t, len(allFields), 6, "Unexpected number of allFields.")         // Pointer fields are ignored by default
	AssertEqualIntAndFailNow(t, len(allFieldNames), 6, "Unexpected number of allFieldNames.") // Pointer fields are ignored by default

	ignoreSettingsOrmSkipped := &IgnoreFieldTypes{OrmSkipped: true, Created: false, Updated: false, Relationals: false}
	allFieldsIgnoreOrmSkipped := GetAllFieldsOfStruct(&person2, ignoreSettingsOrmSkipped)
	allFieldNamesIgnoreOrmSkipped := GetAllFieldNamesOfStruct(&person2, ignoreSettingsOrmSkipped)
	AssertEqualIntAndFailNow(t, len(allFieldsIgnoreOrmSkipped), 5, "Unexpected number of allFieldsIgnoreOrmSkipped.")
	AssertEqualIntAndFailNow(t, len(allFieldNamesIgnoreOrmSkipped), 5, "Unexpected number of allFieldNamesIgnoreOrmSkipped.")

	ignoreSettingsSkipCreated := &IgnoreFieldTypes{OrmSkipped: false, Created: true, Updated: false, Relationals: false}
	allFieldsIgnoreSkipCreated := GetAllFieldsOfStruct(&person2, ignoreSettingsSkipCreated)
	allFieldNamesIgnoreSkipCreated := GetAllFieldNamesOfStruct(&person2, ignoreSettingsSkipCreated)
	AssertEqualIntAndFailNow(t, len(allFieldsIgnoreSkipCreated), 5, "Unexpected number of allFieldsIgnoreSkipCreated.")
	AssertEqualIntAndFailNow(t, len(allFieldNamesIgnoreSkipCreated), 5, "Unexpected number of allFieldNamesIgnoreSkipCreated.")

	ignoreSettingsSkipUpdated := &IgnoreFieldTypes{OrmSkipped: false, Created: false, Updated: true, Relationals: false}
	allFieldsIgnoreSkipUpdated := GetAllFieldsOfStruct(&person2, ignoreSettingsSkipUpdated)
	allFieldNamesIgnoreSkipUpdated := GetAllFieldNamesOfStruct(&person2, ignoreSettingsSkipUpdated)
	AssertEqualIntAndFailNow(t, len(allFieldsIgnoreSkipUpdated), 5, "Unexpected number of allFieldsIgnoreSkipUpdated.")
	AssertEqualIntAndFailNow(t, len(allFieldNamesIgnoreSkipUpdated), 5, "Unexpected number of allFieldNamesIgnoreSkipUpdated.")

	ignoreSettingsSkipRelationals := &IgnoreFieldTypes{OrmSkipped: false, Created: false, Updated: false, Relationals: true}
	allFieldsIgnoreSkipRelationals := GetAllFieldsOfStruct(&person2, ignoreSettingsSkipRelationals)
	allFieldNamesIgnoreSkipRelationals := GetAllFieldNamesOfStruct(&person2, ignoreSettingsSkipRelationals)
	AssertEqualIntAndFailNow(t, len(allFieldsIgnoreSkipRelationals), 6, "Unexpected number of allFieldsIgnoreSkipRelationals.")         // Pointer fields are ignored by default
	AssertEqualIntAndFailNow(t, len(allFieldNamesIgnoreSkipRelationals), 6, "Unexpected number of allFieldNamesIgnoreSkipRelationals.") // Pointer fields are ignored by default

	ignoreSettingsSkipAllFourTypes := &IgnoreFieldTypes{OrmSkipped: true, Created: true, Updated: true, Relationals: true}
	allFieldsIgnoreSkipAllFourTypes := GetAllFieldsOfStruct(&person2, ignoreSettingsSkipAllFourTypes)
	allFieldNamesIgnoreSkipAllFourTypes := GetAllFieldNamesOfStruct(&person2, ignoreSettingsSkipAllFourTypes)
	AssertEqualIntAndFailNow(t, len(allFieldsIgnoreSkipAllFourTypes), 3, "Unexpected number of allFieldsIgnoreSkipAllFourTypes.")
	AssertEqualIntAndFailNow(t, len(allFieldNamesIgnoreSkipAllFourTypes), 3, "Unexpected number of allFieldNamesIgnoreSkipAllFourTypes.")
}

func doTestPassingNonStructPointer_ForTestDoesStructContainField(t *testing.T) {
	defer func() {
		recover()
	}()

	tmpStruct1 := tmpEmptyStruct1{}
	DoesStructContainField(tmpStruct1, "")

	AssertAndFailNow(t, false, "The previous line should not have succeeded. We should require a POINTER, but we are not passing in a pointer.")
}
func TestDoesStructContainField(t *testing.T) {
	doTestPassingNonStructPointer_ForTestDoesStructContainField(t)

	tmpStruct := &struct {
		Name          string
		Age           int
		StructPointer *tmpEmptyStruct1
	}{}

	AssertAndFailNow(t, DoesStructContainField(tmpStruct, "Name"), "Unexpected because the struct does contain field 'Name'")
	AssertAndFailNow(t, DoesStructContainField(tmpStruct, "Age"), "Unexpected because the struct does contain field 'Age'")
	AssertAndFailNow(t, DoesStructContainField(tmpStruct, "StructPointer"), "Unexpected because the struct does contain field 'StructPointer'")
}
