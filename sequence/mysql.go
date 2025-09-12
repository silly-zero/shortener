package sequence

import (
	"database/sql"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

//建立mysql连接, 执行replace into语句
//replace into sequence (stub) values ('a')
//select last_insert_id()

const sqlReplaceIntoStub = "replace into sequence (stub) values ('a')"

type MySQL struct {
	conn sqlx.SqlConn
}

func NewMySQL(dsn string) Sequence {
	return &MySQL{
		conn: sqlx.NewMysql(dsn),
	}
}

// Next 取下一个序号
func (m *MySQL) Next() (seq uint64, err error) {
	//准备语句
	var stmt sqlx.StmtSession
	stmt, err = m.conn.Prepare(sqlReplaceIntoStub)
	if err != nil {
		logx.Errorw("conn.Prepare failed", logx.Field("err", err.Error()))
		return 0, err
	}
	defer stmt.Close()

	//执行
	var rest sql.Result
	rest, err = stmt.Exec()
	if err != nil {
		logx.Errorw("stmt.Exec failed", logx.Field("err", err.Error()))
		return 0, err
	}

	//获取最后插入的id
	var lid int64
	lid, err = rest.LastInsertId()
	if err != nil {
		logx.Errorw("rest.LastInsertId failed", logx.Field("err", err.Error()))
		return 0, err
	}
	return uint64(lid), nil
}
