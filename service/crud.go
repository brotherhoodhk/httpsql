package service

import (
	"errors"
	"httpsql/model"
	"strings"

	"github.com/gin-gonic/gin"
)

var dbparent = model.ROOTPATH + "/database/"
var dbpointer = make(map[string]*model.DataBase) //one database to one pointer
func createdb(ctx *gin.Context) {
	db_name := ctx.Param("database")
	var err error
	code := 200
	data := gin.H{}
	if len(db_name) > 0 && strings.ContainsRune(db_name, ' ') {
		if _, ok := dbpointer[db_name]; !ok {
			dbpointer[db_name] = model.NewDb(db_name)
		} else {
			err = errors.New("database " + db_name + " is existed")
		}
	} else {
		err = errors.New("database name is not format")
	}
	data = gin.H{"msg": err.Error()}
	ctx.JSON(code, data)
}
func create_table(ctx *gin.Context) {
	table_name := ctx.Param("table_name")
	db_name := ctx.Param("database")
	construct := ctx.PostFormMap("construct")
	var code int = 200
	var err error
	data := gin.H{}
	if pointer, ok := dbpointer[db_name]; invalidnameset(table_name, db_name) && ok {
		err = pointer.Create_Table(table_name, construct)
	} else {
		err = errors.New("table name or db name is not correct format")
	}
	data = gin.H{"msg": err.Error()}
	ctx.JSON(code, data)
}
func insertdata(ctx *gin.Context) {
	table_name := ctx.GetHeader("table_name")
	db_name := ctx.GetHeader("database")
	insertdata := ctx.PostFormMap("data")
	var code int = 200
	var err error
	data := gin.H{}
	if pointer, ok := dbpointer[db_name]; invalidnameset(table_name, db_name) && ok {
		err = pointer.InsertData(table_name, insertdata)
	} else {
		err = errors.New("table name or db name is not correct format")
	}
	data = gin.H{"msg": err.Error()}
	ctx.JSON(code, data)
}
func droptable(ctx *gin.Context) {
	table_name := ctx.GetHeader("table_name")
	db_name := ctx.GetHeader("database")
	var code int = 200
	var err error
	data := gin.H{}
	if pointer, ok := dbpointer[db_name]; invalidnameset(table_name, db_name) && ok {
		err = pointer.DropTable(table_name)
	} else {
		err = errors.New("table name or db name is not correct format")
	}
	data = gin.H{"msg": err.Error()}
	ctx.JSON(code, data)
}
func getdata(ctx *gin.Context) {
	table_name := ctx.GetHeader("table_name")
	db_name := ctx.GetHeader("database")
	patternkey := ctx.PostForm("patternkey")
	patternval := ctx.PostForm("patternval")
	args := ctx.PostFormArray("args")
	var code int = 200
	var err error
	var data any
	if pointer, ok := dbpointer[db_name]; invalidnameset(table_name, db_name) && ok {
		resmap, err := pointer.GetData(table_name, patternkey, patternval, args...)
		if err == nil {
			data = resmap
		}
	} else {
		err = errors.New("table name or db name is not correct format")
	}
	if err != nil {
		data = gin.H{"msg": err.Error()}
	}
	ctx.JSON(code, data)
}
