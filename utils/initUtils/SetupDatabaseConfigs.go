package initUtils

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"net"
	"strconv"
	"time"

	. "github.com/francoishill/goangi2/context"
	. "github.com/francoishill/goangi2/utils/configUtils"
	. "github.com/francoishill/goangi2/utils/cookieUtils"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func SetupDatabaseConfigs(configProvider IConfigContainer, ormSyncNow, ormSyncForce, ormSyncIfFlagPresent bool) {
	driverName := configProvider.MustString("database::driver_name")
	dataSource := configProvider.MustString("database::data_source")
	maxIdleConnections := configProvider.DefaultInt("database::max_idle_conn", 30)
	maxOpenConnections := configProvider.DefaultInt("database::max_open_conn", 50)

	orm.DefaultTimeLoc = time.Local
	orm.RegisterDataBase("default", driverName, dataSource, maxIdleConnections, maxOpenConnections)

	if ormSyncNow {
		force := ormSyncForce
		verbose := true
		orm.RunSyncdb("default", force, verbose)
	}

	if ormSyncIfFlagPresent {
		// This will only run if the commandline arguments are "orm ..."
		orm.RunCommand()
	}
}

func SetupDefaultSecurityContext(configProvider IConfigContainer) {
	DefaultCookieSecurityContext = CreateCookieSecurityContext(configProvider.MustString("security::cookie_security_key"))
}

func SetupServerConfigs_AndAppContext(configProvider IConfigContainer) *BaseAppContext {
	baseAppUrl := configProvider.MustString("server::base_app_url")
	hostAndPort := configProvider.MustString("server::host_and_port")
	host, portStr, err := net.SplitHostPort(hostAndPort)
	checkError(err)

	port, err := strconv.ParseInt(portStr, 10, 32)
	checkError(err)

	beego.HttpAddr = host
	beego.HttpPort = int(port)
	beego.HttpServerTimeOut = configProvider.DefaultInt64("server::http_server_timeout", 0)

	uploadDir := configProvider.DefaultString("file_paths::temp_upload_dir", "temp_uploads")
	profilePicsDir := configProvider.DefaultString("file_paths::profile_pics_dir", "profile_pics")

	//Context settings
	maxProfilePicWidth := uint(configProvider.DefaultInt("other::max_profile_pic_width", 128))
	DefaultBaseAppContext = CreateBaseAppContext(beego.BeeLogger, baseAppUrl, maxProfilePicWidth, uploadDir, profilePicsDir)
	return DefaultBaseAppContext
}
