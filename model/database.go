package model

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"time"

	"github.com/oswaldoooo/octools/toolsbox"
)

var datadir = ROOTPATH + "/database/"
var DirMode = 0740
var FileMode = 0640
var dataidlength = 20
var wrbuffersize = 1 << 10

type DataBase struct {
	RootName    string
	table_info  []*Table
	create_date string
}
type Table struct {
	Name      string
	Columnact map[string]*os.File
}

// create the new database pointer
func NewDb(rootname string) (db *DataBase) {
	_, err := os.Stat(datadir + rootname)
	if err != nil {
		err = os.Mkdir(datadir+rootname, os.FileMode(DirMode))
		if err != nil {
			//record the error information
			Errorlog.Println(err.Error())
		}
	}
	return &DataBase{RootName: rootname, create_date: time.Now().Format(time.ANSIC)}
}

// construct filedname=filedtype,column format is name****type
func (s *DataBase) Create_Table(table_name string, construct map[string]string) (err error) {
	_, err = os.Stat(datadir + s.RootName + "/" + table_name)
	if err != nil {
		parentpath := datadir + s.RootName + "/" + table_name + "/"
		err = os.Mkdir(datadir+s.RootName+"/"+table_name, fs.ModeDir)
		if err == nil {
			tabler := &Table{Name: table_name, Columnact: make(map[string]*os.File)}
			//create the column.format is name****type
			for name, typeinfo := range construct {
				act, err := os.OpenFile(parentpath+name+"****"+typeinfo+".db", os.O_CREATE|os.O_WRONLY, os.FileMode(FileMode))
				if err != nil {
					Errorlog.Println(err)
					break
				} else {
					//put column into tabler
					tabler.Columnact[name] = act
				}
			}
			if err == nil {
				//link tabler to db
				s.table_info = append(s.table_info, tabler)
			}
		}
	} else {
		err = errors.New("create table " + table_name + " failed it's existed")
	}
	return
}

// insertdata map is [column name : value]
func (s *DataBase) InsertData(table_name string, insertdata map[string]string) (err error) {
	var tableree *Table
	for _, tabler := range s.table_info {
		if tabler.Name == table_name {
			tableree = tabler
			break
		}
	}
	if tableree == nil {
		err = errors.New(table_name + " dont existed")
		return
	}
	keysarr := toolsbox.ExportMapKeys(insertdata)
	if toolsbox.CheckArgs(keysarr, tableree.Columnact) {
		//compare the insert data's columns and origin table columns
		dataid := generateid(dataidlength)
		for name, value := range insertdata {
			actor := tableree.Columnact[name]
			_, err = actor.Write([]byte(dataid + "****" + value))
			if err != nil {
				Errorlog.Println(err.Error())
				break
			}
		}
	} else {
		err = errors.New("columns are not correct ")
	}
	return
}

// get the data array
func (s *DataBase) GetData(table_name, patternkey, patternvalue string, columns ...string) (res []map[string]string, err error) {
	var tableree *Table
	for _, tabler := range s.table_info {
		if tabler.Name == table_name {
			tableree = tabler
			break
		}
	}
	if tableree == nil {
		err = errors.New(table_name + " dont existed")
		return
	}
	res = []map[string]string{}
	var idarr []string
	//if the column existed get the id
	if act, ok := tableree.Columnact[patternkey]; ok {
		buffers := make([]byte, wrbuffersize)
		read := bufio.NewReader(act)
		lang, _ := read.Read(buffers)
		//get the dataarray's id.
		idarr, err = ParseDbContent(buffers[:lang])
		if err != nil {
			return
		} else if len(idarr) < 1 {
			err = errors.New("id is dont existed")
			return
		}
	} else {
		//pattern column dont existed
		err = errors.New(patternkey + " dont existed")
		return
	}
	if len(columns) == 0 {

	} else {
		//from columns get the target value
		if toolsbox.CheckArgs(columns, tableree.Columnact) {
			for _, id := range idarr {
				idresmap := make(map[string]string)
				for _, name := range columns {
					resmap, err := ParseDb(tableree.Columnact[name])
					if err == nil {
						if value, ok := resmap[id]; ok {
							//get target value
							idresmap[name] = value
						}
					} else {
						//error occured
						break
					}
				}
				if err == nil {
					//if error dont happend ,put columns into res array
					res = append(res, idresmap)
				}
			}
		} else {
			err = errors.New("columns dont existed")
		}
	}
	return
}

// delete the target data
func (s *DataBase) DeleteData(table_name, patternkey, patternvalue string) (err error) {
	var tableree *Table
	for _, tabler := range s.table_info {
		if tabler.Name == table_name {
			tableree = tabler
			break
		}
	}
	if tableree == nil {
		err = errors.New(table_name + " dont existed")
		return
	}
	var idarr []string
	//if the column existed get the id
	if act, ok := tableree.Columnact[patternkey]; ok {
		buffers := make([]byte, wrbuffersize)
		read := bufio.NewReader(act)
		lang, _ := read.Read(buffers)
		//get the dataarray's id.
		idarr, err = ParseDbContent(buffers[:lang])
		if err != nil {
			return
		} else if len(idarr) < 1 {
			err = errors.New("id is dont existed")
			return
		}
	} else {
		//pattern column dont existed
		err = errors.New(patternkey + " dont existed")
		return
	}
	for _, columns := range tableree.Columnact {
		deletetargetvale(idarr, columns)
	}
	return
}

// drop the target table
func (s *DataBase) DropTable(table_name string) (err error) {
	//--------make sure the table is existed--------
	var tableree *Table
	for ke, tabler := range s.table_info {
		if tabler.Name == table_name {
			tableree = tabler
			if ke == len(s.table_info)-1 {
				s.table_info = s.table_info[:ke]
			} else {
				s.table_info = append(s.table_info[:ke], s.table_info[ke+1:]...)
			}
			break
		}
	}
	if tableree == nil {
		err = errors.New(table_name + " dont existed")
		return
	}
	for _, actor := range tableree.Columnact {
		//close the file open
		actor.Close()
	}
	//delete the table dir
	err = os.RemoveAll(datadir + s.RootName + "/" + table_name)
	return
}
