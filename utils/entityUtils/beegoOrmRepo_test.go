package entityUtils

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

const (
	cTEMP_BEEGO_ORM_REPO_TEST__SQLITE_FILE = "beego_orm_repo_test.sqlite"
)

type Post struct {
	Id    int64 `orm:"pk;auto"`
	Title string
	Tags  []*Tag `orm:"rel(m2m)"`
}

type Tag struct {
	Id    int64 `orm:"pk;auto"`
	Name  string
	Posts []*Post `orm:"reverse(many)"`
}

func TestBaseListEntities_OrderBy_Limit_Offset__WhereM2MCountIsZero(t *testing.T) {
	Convey("Testing BaseListEntities_ANDFilters_OrderBy_Limit_Offset", t, func() {
		orm.RegisterModel(new(Post), new(Tag))

		orm.DefaultTimeLoc = time.Local
		maxIdleConnections := 30
		maxOpenConnections := 50
		err := orm.RegisterDataBase("default", "sqlite3", cTEMP_BEEGO_ORM_REPO_TEST__SQLITE_FILE, maxIdleConnections, maxOpenConnections)
		if err != nil {
			panic(err)
		}

		force := true
		verbose := true
		err = orm.RunSyncdb("default", force, verbose)
		if err != nil {
			panic(err)
		}

		ormContext := CreateDefaultOrmContext()

		post1 := &Post{Title: "Post1"}
		OrmRepo.BaseInsertEntity(ormContext, post1)

		post2 := &Post{Title: "Post2"}
		OrmRepo.BaseInsertEntity(ormContext, post2)
		tag1 := &Tag{Name: "Tag1"}
		OrmRepo.BaseInsertEntity(ormContext, tag1)
		tag2 := &Tag{Name: "Tag1"}
		OrmRepo.BaseInsertEntity(ormContext, tag2)
		OrmRepo.BaseUpdateM2MByAddAndRemove(ormContext, post2, "tags", []interface{}{}, []interface{}{tag1, tag2})

		post3 := &Post{Title: "Post3"}
		OrmRepo.BaseInsertEntity(ormContext, post3)

		/*tags := []*Tag{}
		OrmRepo.BaseListEntities_OrderBy_Limit_Offset__WhereM2MCountIsZero(ormContext, "post", []string{}, 1000, 0, nil, tags)
		Convey("Expect the 'BaseListEntities_OrderBy_Limit_Offset__WhereM2MCountIsZero' to return post 2", func() {
			So(len(tags), ShouldEqual, 2)
			So(tags[0].Name, ShouldEqual, "Post1")
			So(tags[2].Name, ShouldEqual, "Post3")
		})
		*/
	})
	//This does not work, probably due to a lock: os.Remove(cTEMP_BEEGO_ORM_REPO_TEST__SQLITE_FILE)
}
