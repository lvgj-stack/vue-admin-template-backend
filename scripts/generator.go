package main

import (
	"context"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"github.com/Mr-LvGJ/gobase/pkg/common/setting"
	"github.com/Mr-LvGJ/gobase/pkg/gobase/store"
)

// generate code
func main() {
	setting.InitConfig("../configs/dev.gobase.yaml")
	DB, err := store.Setup(context.Background())
	if err != nil {
		panic(err)
	}
	// specify the output directory (default: "./query")
	// ### if you want to query without context constrain, set mode gen.WithoutContext ###
	g := gen.NewGenerator(gen.Config{
		Mode:         gen.WithDefaultQuery,
		OutPath:      "../pkg/gobase/store/dal",
		ModelPkgPath: "../pkg/gobase/model/v1/entity",
		/* Mode: gen.WithoutContext,*/
		// if you want the nullable field generation property to be pointer type, set FieldNullable true
		FieldNullable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})

	// reuse the database connection in Project or create a connection here
	// if you want to use GenerateModel/GenerateModelAs, UseDB is necessary, or it will panic
	g.UseDB(DB)

	updatedAtTagF := gen.FieldGORMTag("updated_at", func(tag field.GormTag) field.GormTag {
		tag.Set("default", "current_timestamp(6) on update current_timestamp(6)")
		return tag
	})

	// apply basic crud api on structs or table models which is specified by table name with function
	// GenerateModel/GenerateModelAs. And generator will generate table models' code when calling Excute.
	g.ApplyBasic(
		g.GenerateModel("user", updatedAtTagF),
		g.GenerateModel("table", updatedAtTagF, gen.FieldType("status", "constant.TableStatus")),
		//g.GenerateModel("attachment", gen.FieldType("type", "consts.AttachmentType")),
		//g.GenerateModel("category", gen.FieldType("type", "consts.CategoryType")),
		// g.GenerateModel("comment", gen.FieldType("type", "consts.CommentType"), gen.FieldType("status",
		// "consts.CommentStatus")),
		//g.GenerateModel("comment_black"),
		//g.GenerateModel("flyway_schema_history"),
		//g.GenerateModel("journal", gen.FieldType("type", "consts.JournalType")),
		//g.GenerateModel("link"),
		//g.GenerateModel("log", gen.FieldType("type", "consts.LogType")),
		//g.GenerateModel("menu"),
		//g.GenerateModelAs("meta", "Meta", gen.FieldType("type", "consts.MetaType")),
		//g.GenerateModel("option", gen.FieldType("type", "consts.OptionType")),
		//g.GenerateModel("photo"),
		// g.GenerateModel("post", gen.FieldType("type", "consts.PostType"), gen.FieldType("status",
		// "consts.PostStatus"), gen.FieldType("editor_type", "consts.EditorType")),
		//g.GenerateModel("post_category"),
		//g.GenerateModel("post_tag"),
		//g.GenerateModel("tag"),
		//g.GenerateModel("theme_setting"),
		//g.GenerateModel("user", gen.FieldType("mfa_type", "consts.MFAType")),
		//g.GenerateModel("temporary"),
		//g.GenerateModel("permanent"),
	)

	// apply diy interfaces on structs or table models
	// g.ApplyInterface(func(method model.Method) {}, model.User{}, g.GenerateModel("company"))

	// execute the action of code generation
	g.Execute()
}
