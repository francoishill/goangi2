package translationUtils

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetTranslationWithParams(t *testing.T) {
	Convey("GetTranslationWithParams", t, func() {
		var map1 = `{
			"KEY_OUTER_1": {
				"KEY_INNER_1": {
					"KEY1": "My value 1"
				}
			}
		}`
		var map2 = `{
			"KEY_OUTER_2": {
				"KEY_INNER_2": {
					"KEY2": "My value 2 with {{personName}} and I am {{age}} years old of age"
				}
			}
		}`
		var map3 = `{
			"KEY_OUTER_3": {
				"KEY_INNER_3": {
					"KEY3": "My value 3"
				}
			}
		}`

		Convey("Check that our value is correctly translated", func() {
			provider := CreateTranslationProvider([]*translationsMap{
				CreateTranslationsMap(map1),
				CreateTranslationsMap(map2),
				CreateTranslationsMap(map3),
			})

			params := map[string]string{
				"personName": "My Full Name",
				"age":        "25",
			}

			val2 := provider.GetTranslationWithParams("KEY_OUTER_2.KEY_INNER_2.KEY2", params)
			So(val2, ShouldEqual, "My value 2 with My Full Name and I am 25 years old of age")

			val1 := provider.GetTranslation("KEY_OUTER_1.KEY_INNER_1.KEY1")
			So(val1, ShouldEqual, "My value 1")

			val3 := provider.GetTranslation("KEY_OUTER_3.KEY_INNER_3.KEY3")
			So(val3, ShouldEqual, "My value 3")

			valInvalid := provider.GetTranslation("KEY.INVALID")
			So(valInvalid, ShouldEqual, "KEY.INVALID{}")

			valInvalid = provider.GetTranslationWithParams("KEY.INVALID", map[string]string{"ParamKey1": "ParamValue1", "PK2": "PV2"})
			So(valInvalid, ShouldEqual, `KEY.INVALID{"ParamKey1":"ParamValue1", "PK2":"PV2"}`)
		})
	})
}
