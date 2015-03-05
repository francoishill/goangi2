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
	Convey("doTestPassingNonStructPointer_ForTestSetFormValues", t, func() {
		So(func() { SetFormValues(obj1, obj2, nil, nil) }, ShouldPanic)
	})
}
func TestSetFormValues(t *testing.T) {
	tmpStruct1 := tmpPersonStructForSetFormValues{}
	tmpStruct2 := tmpPersonStructForSetFormValues{}
	doTestPassingNonStructPointer_ForTestSetFormValues(t, tmpStruct1, tmpStruct2, "Case 1")  //Both non pointers
	doTestPassingNonStructPointer_ForTestSetFormValues(t, &tmpStruct1, tmpStruct2, "Case 2") //Only one not pointer
	doTestPassingNonStructPointer_ForTestSetFormValues(t, tmpStruct1, &tmpStruct2, "Case 3") //Only one not pointer

	var destinationPersonToUpdate *tmpPersonStructForSetFormValues

	Convey("Having a source and destination person (case 1)", t, func() {
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
		Convey("Unexpected Name (case 1) for destinationPersonToUpdate.", func() {
			So(destinationPersonToUpdate.Name, ShouldEqual, "OrigName")
		})
		Convey("Unexpected Age (case 1) for destinationPersonToUpdate.", func() {
			So(destinationPersonToUpdate.Age, ShouldEqual, 25)
		})
		Convey("Will now SetFormValues", func() {
			SetFormValues(sourcePersonToCopyFrom, destinationPersonToUpdate, nil, nil) //Transfer from sourcePersonToCopyFrom to destinationPersonToUpdate
			Convey("Name (case 1) did not change with SetFormValues.", func() {
				So(destinationPersonToUpdate.Name, ShouldEqual, "NewName")
			})
			Convey("Age (case 1) did not change with SetFormValues.", func() {
				So(destinationPersonToUpdate.Age, ShouldEqual, 26)
			})
		})
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
	Convey("Having a source and destination person (case 2)", t, func() {
		Convey("Unexpected Name (case 2) for destinationPersonToUpdate.", func() {
			So(destinationPersonToUpdate.Name, ShouldEqual, "OrigName")
		})
		Convey("Unexpected Age (case 2) for destinationPersonToUpdate.", func() {
			So(destinationPersonToUpdate.Age, ShouldEqual, 25)
		})
		Convey("Will now SetFormValues", func() {
			SetFormValues(sourceAnonToCopyFrom, destinationPersonToUpdate, nil, nil)
			Convey("Name (case 2) did not change with SetFormValues.", func() {
				So(destinationPersonToUpdate.Name, ShouldEqual, "NewNameAnon")
			})
			Convey("Age (case 2) did not change with SetFormValues.", func() {
				So(destinationPersonToUpdate.Age, ShouldEqual, 27)
			})
		})
	})

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
	Convey("Having a source and destination person (case 3)", t, func() {
		Convey("Unexpected Name (case 3) for destinationPersonToUpdate.", func() {
			So(destinationPersonToUpdate.Name, ShouldEqual, "OrigName")
		})
		Convey("Unexpected Age (case 3) for destinationPersonToUpdate.", func() {
			So(destinationPersonToUpdate.Age, ShouldEqual, 25)
		})
		Convey("Will now SetFormValues and skip Name and Age", func() {
			SetFormValues(sourceAnonToCopyFrom2, destinationPersonToUpdate, []string{"Name", "Age"}, nil)
			Convey("Name (case 3) was NOT supposed to update with SetFormValues as we ignored Name+Age.", func() {
				So(destinationPersonToUpdate.Name, ShouldEqual, "OrigName")
			})
			Convey("Age (case 3) was NOT supposed to update with SetFormValues as we ignored Name+Age.", func() {
				So(destinationPersonToUpdate.Age, ShouldEqual, 25)
			})

			Convey("Will now SetFormValues and skip Name only", func() {
				SetFormValues(sourceAnonToCopyFrom2, destinationPersonToUpdate, []string{"Name"}, nil)
				Convey("Name (case 4) was NOT supposed to update with SetFormValues as we ignored Name.", func() {
					So(destinationPersonToUpdate.Name, ShouldEqual, "OrigName")
				})
				Convey("Age (case 4) did not change with SetFormValues.", func() {
					So(destinationPersonToUpdate.Age, ShouldEqual, 27)
				})
			})
		})
	})

	//Just reset the values of Name+Age
	destinationPersonToUpdate = &tmpPersonStructForSetFormValues{
		Name:     "OrigName",
		Age:      25,
		Employed: true,
	}
	Convey("Having a source and destination person (case 4)", t, func() {
		Convey("Will now SetFormValues and skip Age", func() {
			SetFormValues(sourceAnonToCopyFrom2, destinationPersonToUpdate, []string{"Age"}, nil)
			So(destinationPersonToUpdate.Name, ShouldEqual, "NewNameAnon")
			So(destinationPersonToUpdate.Age, ShouldEqual, 25)
		})
	})
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
	Convey("Getting FormChanges (case 1)", t, func() {
		changedFields1 := FormChanges(originalPerson, modifiedPerson1, nil, nil)
		Convey("Unexpected number of changed fields for identical entity.", func() {
			So(len(changedFields1), ShouldEqual, 0)
		})
	})

	//Only modify Name
	modifiedPerson2 := &tmpPersonStructForFormChanges{
		Name:           "NewName",
		SurnameIgnored: "Surname",
		Age:            24,
		Employed:       true,
		Created:        now,
		Updated:        now,
	}
	Convey("Getting FormChanges (case 2)", t, func() {
		changedFields2 := FormChanges(originalPerson, modifiedPerson2, nil, nil)
		Convey("Unexpected number of changed fields for only modifying only the Name field.", func() {
			So(len(changedFields2), ShouldEqual, 1)
		})

		changedField2_0 := changedFields2[0]
		Convey("Unexpected changed FieldName.", func() {
			So(changedField2_0.FieldName, ShouldEqual, "Name")
		})
		Convey("Unexpected changed OldValue.", func() {
			So(changedField2_0.OldValue, ShouldEqual, "Name")
		})
		Convey("Unexpected changed NewValue.", func() {
			So(changedField2_0.NewValue, ShouldEqual, "NewName")
		})
	})

	//Only modify Age
	modifiedPerson3 := &tmpPersonStructForFormChanges{
		Name:           "Name",
		SurnameIgnored: "Surname",
		Age:            25,
		Employed:       true,
		Created:        now,
		Updated:        now,
	}
	Convey("Getting FormChanges (case 3)", t, func() {
		changedFields3 := FormChanges(originalPerson, modifiedPerson3, nil, nil)
		Convey("Unexpected number of changed fields for only modifying only the Age field.", func() {
			So(len(changedFields3), ShouldEqual, 1)
		})
		changedField3_0 := changedFields3[0]
		Convey("Unexpected changed FieldName.", func() {
			So(changedField3_0.FieldName, ShouldEqual, "Age")
		})
		Convey("Unexpected changed OldValue.", func() {
			So(changedField3_0.OldValue, ShouldEqual, "24")
		})
		Convey("Unexpected changed NewValue.", func() {
			So(changedField3_0.NewValue, ShouldEqual, "25")
		})
	})

	//Only modify Employed
	modifiedPerson4 := &tmpPersonStructForFormChanges{
		Name:           "Name",
		SurnameIgnored: "Surname",
		Age:            24,
		Employed:       false,
		Created:        now,
		Updated:        now,
	}
	Convey("Getting FormChanges (case 4)", t, func() {
		changedFields4 := FormChanges(originalPerson, modifiedPerson4, nil, nil)
		Convey("Unexpected number of changed fields for only modifying only the Employed field.", func() {
			So(len(changedFields4), ShouldEqual, 1)
		})
		changedField4_0 := changedFields4[0]
		Convey("Unexpected changed FieldName.", func() {
			So(changedField4_0.FieldName, ShouldEqual, "Employed")
		})
		Convey("Unexpected changed OldValue.", func() {
			So(changedField4_0.OldValue, ShouldEqual, "true")
		})
		Convey("Unexpected changed NewValue.", func() {
			So(changedField4_0.NewValue, ShouldEqual, "false")
		})
	})

	//Modify all fields
	modifiedPerson5 := &tmpPersonStructForFormChanges{
		Name:           "NewName",
		SurnameIgnored: "NewSurname",
		Age:            26,
		Employed:       false,
		Created:        now,
		Updated:        now,
	}
	Convey("Getting FormChanges (case 5)", t, func() {
		changedFields5a := FormChanges(originalPerson, modifiedPerson5, nil, nil)
		Convey("Unexpected number of changed fields for modifying Name, Age and Employed (Surname ignored) fields.", func() {
			So(len(changedFields5a), ShouldEqual, 3)
		})
		changedFields5b := FormChanges(originalPerson, modifiedPerson5, []string{"Name", "Age"}, nil) //Skipping the changes of Name+Age changes
		Convey("Unexpected number of changed fields for modifying Name, Age and Employed (Surname ignored and skipping Name+Age) fields.", func() {
			So(len(changedFields5b), ShouldEqual, 1)
		})
	})

	//Modify nothing but only the Created and Updated, but these two fields should be ignored.
	modifiedPerson6 := &tmpPersonStructForFormChanges{
		Name:           "Name",
		SurnameIgnored: "Surname",
		Age:            24,
		Employed:       true,
		Created:        now.Add(10 * time.Hour),
		Updated:        now.Add(10 * time.Hour),
	}
	Convey("Getting FormChanges (case 5)", t, func() {
		changedFields6 := FormChanges(originalPerson, modifiedPerson6, nil, nil)
		Convey("Unexpected number of changed fields for entity where ONLY the Created and Updated dates changed. These dates are expected to be ignored.", func() {
			So(len(changedFields6), ShouldEqual, 0)
		})
	})

	//Modify nothing but only the form ignored field
	modifiedPerson7 := &tmpPersonStructForFormChanges{
		Name:           "Name",
		SurnameIgnored: "NewSurname",
		Age:            24,
		Employed:       true,
		Created:        now,
		Updated:        now,
	}
	Convey("Getting FormChanges (case 6)", t, func() {
		changedFields7 := FormChanges(originalPerson, modifiedPerson7, nil, nil)
		Convey("Unexpected number of changed fields for entity where ONLY the SurnameIgnored changed, which is a field to be ignored..", func() {
			So(len(changedFields7), ShouldEqual, 0)
		})
	})
}

type tmpEmptyStruct1 struct{}
type tmpEmptyStruct2 struct{}

func doTestPassingNonStructPointer_ForTestGetAllFieldsOfStruct(t *testing.T) {
	Convey("This call should panic. We should require a POINTER, but we are not passing in a pointer.", t, func() {
		tmpStruct1 := tmpEmptyStruct1{}
		ignoreSettings := &IgnoreFieldTypes{OrmSkipped: false, Created: false, Updated: false, Relationals: false}
		So(func() { GetAllFieldsOfStruct(tmpStruct1, ignoreSettings) }, ShouldPanic)
	})
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
	Convey("TestGetAllFieldsOfStruct_AndGetAllFieldNamesOfStruct", t, func() {
		ignoreSettingsNone := &IgnoreFieldTypes{OrmSkipped: false, Created: false, Updated: false, Relationals: false}
		allFields := GetAllFieldsOfStruct(&person2, ignoreSettingsNone)
		allFieldNames := GetAllFieldNamesOfStruct(&person2, ignoreSettingsNone)
		Convey("Unexpected number of allFields (pointer fields are ignored by default).", func() {
			So(len(allFields), ShouldEqual, 6)
		})
		Convey("Unexpected number of allFieldNames (pointer fields are ignored by default).", func() {
			So(len(allFieldNames), ShouldEqual, 6)
		})

		ignoreSettingsOrmSkipped := &IgnoreFieldTypes{OrmSkipped: true, Created: false, Updated: false, Relationals: false}
		allFieldsIgnoreOrmSkipped := GetAllFieldsOfStruct(&person2, ignoreSettingsOrmSkipped)
		allFieldNamesIgnoreOrmSkipped := GetAllFieldNamesOfStruct(&person2, ignoreSettingsOrmSkipped)
		Convey("Unexpected number of allFieldsIgnoreOrmSkipped.", func() {
			So(len(allFieldsIgnoreOrmSkipped), ShouldEqual, 5)
		})
		Convey("Unexpected number of allFieldNamesIgnoreOrmSkipped.", func() {
			So(len(allFieldNamesIgnoreOrmSkipped), ShouldEqual, 5)
		})

		ignoreSettingsSkipCreated := &IgnoreFieldTypes{OrmSkipped: false, Created: true, Updated: false, Relationals: false}
		allFieldsIgnoreSkipCreated := GetAllFieldsOfStruct(&person2, ignoreSettingsSkipCreated)
		allFieldNamesIgnoreSkipCreated := GetAllFieldNamesOfStruct(&person2, ignoreSettingsSkipCreated)
		Convey("Unexpected number of allFieldsIgnoreSkipCreated.", func() {
			So(len(allFieldsIgnoreSkipCreated), ShouldEqual, 5)
		})
		Convey("Unexpected number of allFieldNamesIgnoreSkipCreated.", func() {
			So(len(allFieldNamesIgnoreSkipCreated), ShouldEqual, 5)
		})

		ignoreSettingsSkipUpdated := &IgnoreFieldTypes{OrmSkipped: false, Created: false, Updated: true, Relationals: false}
		allFieldsIgnoreSkipUpdated := GetAllFieldsOfStruct(&person2, ignoreSettingsSkipUpdated)
		allFieldNamesIgnoreSkipUpdated := GetAllFieldNamesOfStruct(&person2, ignoreSettingsSkipUpdated)
		Convey("Unexpected number of allFieldsIgnoreSkipUpdated.", func() {
			So(len(allFieldsIgnoreSkipUpdated), ShouldEqual, 5)
		})
		Convey("Unexpected number of allFieldNamesIgnoreSkipUpdated.", func() {
			So(len(allFieldNamesIgnoreSkipUpdated), ShouldEqual, 5)
		})

		ignoreSettingsSkipRelationals := &IgnoreFieldTypes{OrmSkipped: false, Created: false, Updated: false, Relationals: true}
		allFieldsIgnoreSkipRelationals := GetAllFieldsOfStruct(&person2, ignoreSettingsSkipRelationals)
		allFieldNamesIgnoreSkipRelationals := GetAllFieldNamesOfStruct(&person2, ignoreSettingsSkipRelationals)
		Convey("Unexpected number of allFieldsIgnoreSkipRelationals (pointer fields are ignored by default).", func() {
			So(len(allFieldsIgnoreSkipRelationals), ShouldEqual, 6)
		})
		Convey("Unexpected number of allFieldNamesIgnoreSkipRelationals (pointer fields are ignored by default).", func() {
			So(len(allFieldNamesIgnoreSkipRelationals), ShouldEqual, 6)
		})

		ignoreSettingsSkipAllFourTypes := &IgnoreFieldTypes{OrmSkipped: true, Created: true, Updated: true, Relationals: true}
		allFieldsIgnoreSkipAllFourTypes := GetAllFieldsOfStruct(&person2, ignoreSettingsSkipAllFourTypes)
		allFieldNamesIgnoreSkipAllFourTypes := GetAllFieldNamesOfStruct(&person2, ignoreSettingsSkipAllFourTypes)
		Convey("Unexpected number of allFieldsIgnoreSkipAllFourTypes.", func() {
			So(len(allFieldsIgnoreSkipAllFourTypes), ShouldEqual, 3)
		})
		Convey("Unexpected number of allFieldNamesIgnoreSkipAllFourTypes.", func() {
			So(len(allFieldNamesIgnoreSkipAllFourTypes), ShouldEqual, 3)
		})
	})
}

func doTestPassingNonStructPointer_ForTestDoesStructContainField(t *testing.T) {
	Convey("This call should panic. We should require a POINTER, but we are not passing in a pointer.", t, func() {
		tmpStruct1 := tmpEmptyStruct1{}
		So(func() { DoesStructContainField(tmpStruct1, "") }, ShouldPanic)
	})
}
func TestDoesStructContainField(t *testing.T) {
	doTestPassingNonStructPointer_ForTestDoesStructContainField(t)

	tmpStruct := &struct {
		Name          string
		Age           int
		StructPointer *tmpEmptyStruct1
	}{}

	Convey("TestDoesStructContainField", t, func() {
		Convey("Unexpected because the struct does contain field 'Name'", func() {
			So(DoesStructContainField(tmpStruct, "Name"), ShouldBeTrue)
		})
		Convey("Unexpected because the struct does contain field 'Age'", func() {
			So(DoesStructContainField(tmpStruct, "Age"), ShouldBeTrue)
		})
		Convey("Unexpected because the struct does contain field 'StructPointer'", func() {
			So(DoesStructContainField(tmpStruct, "StructPointer"), ShouldBeTrue)
		})
	})
}
