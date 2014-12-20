package config

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
)

func CreateNewBeegoConfigContainer() *BeegoConfigContainer {
	//Ensure we can parse the config file
	err := beego.ParseConfig()
	if err != nil {
		panic(err)
	}
	//Now return AppConfig of beego
	return &BeegoConfigContainer{
		beegoContainer: beego.AppConfig,
	}
}

type BeegoConfigContainer struct {
	beegoContainer config.ConfigContainer
}

func (this *BeegoConfigContainer) Set(key, val string) error {
	return this.beegoContainer.Set(key, val)
}
func (this *BeegoConfigContainer) String(key string) string {
	return this.beegoContainer.String(key)
}
func (this *BeegoConfigContainer) Strings(key string) []string {
	return this.beegoContainer.Strings(key)
}
func (this *BeegoConfigContainer) Int(key string) (int, error) {
	return this.beegoContainer.Int(key)
}
func (this *BeegoConfigContainer) Int64(key string) (int64, error) {
	return this.beegoContainer.Int64(key)
}
func (this *BeegoConfigContainer) Bool(key string) (bool, error) {
	return this.beegoContainer.Bool(key)
}
func (this *BeegoConfigContainer) Float(key string) (float64, error) {
	return this.beegoContainer.Float(key)
}

func (this *BeegoConfigContainer) DefaultString(key string, defaultval string) string {
	return this.beegoContainer.DefaultString(key, defaultval)
}
func (this *BeegoConfigContainer) DefaultStrings(key string, defaultval []string) []string {
	return this.beegoContainer.DefaultStrings(key, defaultval)
}
func (this *BeegoConfigContainer) DefaultInt(key string, defaultval int) int {
	return this.beegoContainer.DefaultInt(key, defaultval)
}
func (this *BeegoConfigContainer) DefaultInt64(key string, defaultval int64) int64 {
	return this.beegoContainer.DefaultInt64(key, defaultval)
}
func (this *BeegoConfigContainer) DefaultBool(key string, defaultval bool) bool {
	return this.beegoContainer.DefaultBool(key, defaultval)
}
func (this *BeegoConfigContainer) DefaultFloat(key string, defaultval float64) float64 {
	return this.beegoContainer.DefaultFloat(key, defaultval)
}

func (this *BeegoConfigContainer) MustString(key string) string {
	const cDEFAULT_VAL = "NO_VALUE_SPECIFIED"
	val := this.beegoContainer.DefaultString(key, cDEFAULT_VAL)
	if val == cDEFAULT_VAL {
		panic("Config String not found with key: " + key)
	}
	return val
}
func (this *BeegoConfigContainer) MustStrings(key string) []string {
	cDEFAULT_VAL := []string{"NO_VALUE_SPECIFIED"}
	val := this.beegoContainer.DefaultStrings(key, cDEFAULT_VAL)
	if len(val) == len(cDEFAULT_VAL) && val[0] == cDEFAULT_VAL[0] {
		panic("Config Strings not found with key: " + key)
	}
	return val
}
func (this *BeegoConfigContainer) MustInt(key string) int {
	const cDEFAULT_VAL = -9878978
	val := this.beegoContainer.DefaultInt(key, cDEFAULT_VAL)
	if val == cDEFAULT_VAL {
		panic("Config Int not found with key: " + key)
	}
	return val
}
func (this *BeegoConfigContainer) MustInt64(key string) int64 {
	const cDEFAULT_VAL = -9988997789
	val := this.beegoContainer.DefaultInt64(key, cDEFAULT_VAL)
	if val == cDEFAULT_VAL {
		panic("Config Int64 not found with key: " + key)
	}
	return val
}
func (this *BeegoConfigContainer) MustBool(key string) bool {
	testValTrue := this.beegoContainer.DefaultBool(key, true)
	testValFalse := this.beegoContainer.DefaultBool(key, false)
	if testValTrue == true &&
		testValFalse == false {
		panic("Config Bool not found with key: " + key)
	}
	return testValFalse
}

func (this *BeegoConfigContainer) DIY(key string) (interface{}, error) {
	return this.beegoContainer.DIY(key)
}
func (this *BeegoConfigContainer) GetSection(section string) (map[string]string, error) {
	return this.beegoContainer.GetSection(section)
}
func (this *BeegoConfigContainer) SaveConfigFile(filename string) error {
	return this.beegoContainer.SaveConfigFile(filename)
}
