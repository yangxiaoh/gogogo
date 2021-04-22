package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	name string
)

// 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？
// 应当用fmt.Errorf包装错误返回，这样可以给调用者提供更多信息
func main() {
	e := getNameById(&name, 3)

	if e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			// do something
			fmt.Println("没有结果")
		} else {
			// error
			fmt.Printf("err:[%+v]", e)
		}
	}

	fmt.Println(name)
}

func getNameById(name *string, id int) error {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/xy")

	if err != nil {
		return fmt.Errorf("conn err: %w", err)
	}

	defer db.Close()

	err = db.QueryRow("select User from mytable where id=?", id).Scan(name)

	if err != nil {
		return fmt.Errorf("query err: %w", err)
	}

	return nil
}
