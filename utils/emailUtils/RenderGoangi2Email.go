package emailUtils

import (
	"bytes"
	"github.com/astaxie/beego"
	"io/ioutil"
)

func RenderGoangi2Email(templatePath string, emailDataObject map[interface{}]interface{}) string {
	if beego.RunMode == "dev" {
		beego.BuildTemplate(beego.ViewsPath)
	}

	if _, ok := beego.BeeTemplates[templatePath]; !ok {
		panic("Unable to find beego template file in path:" + templatePath)
	}

	ibytes := bytes.NewBufferString("")
	err := beego.BeeTemplates[templatePath].ExecuteTemplate(ibytes, templatePath, emailDataObject)
	if err != nil {
		panic("Error executing template: " + err.Error())
	}
	icontent, err := ioutil.ReadAll(ibytes)
	if err != nil {
		panic(err)
	}
	return string(icontent)
}
