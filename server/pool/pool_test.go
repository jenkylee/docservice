package pool

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io"
	"testing"
	"time"
)

func TestNewCommonPool(t *testing.T) {
	pool, err := NewCommonPool(10, 100, time.Second*300, poolMysqlFactory)
	if err != nil {
		t.Fatal("数据库连接错误")
	}
	db, err := pool.Acquire();
	if err != nil {
		t.Fatal("未获取到资源")
	}

	err = pool.Release(db)
	if err != nil {
		t.Fatal("释放资源失败")
	}
	err = pool.Shutdown()
	if err != nil {
		t.Fatal("关闭资源池失败")
	}

	fmt.Println("执行完成")
}

// mysql pool factory
func poolMysqlFactory() (io.Closer, error) {
	db, err := gorm.Open("mysql", "root:123qwe@/mysql?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}

	return db, nil
}
