package example

import (
	"fmt"

	gobatis "github.com/runner-mei/GoBatis"
	"github.com/runner-mei/GoBatis/tests"
)

func ExampleSimple() {
	insertUser := AuthUser{
		Username: "abc",
		Phone:    "123",
		Status:   1,
	}

	gobatis.ShowSQL = false

	factory, err := gobatis.New(&gobatis.Config{DriverName: tests.TestDrv,
		DataSource: tests.TestConnURL,
		//XMLPaths: []string{"example/test.xml"},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		if err = factory.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	switch tests.TestDrv {
	case "postgres":
		_, err = factory.DB().Exec(postgres)
	default:
		_, err = factory.DB().Exec(mysql)
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	ref := factory.Reference()
	userDao := NewAuthUserDao(&ref)
	id, err := userDao.Insert(&insertUser)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("insert success!")

	u, err := userDao.Get(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("fetch user from database!")
	fmt.Println(u.Username)

	_, err = userDao.Delete(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("delete success!")

	tx, err := factory.Begin()
	if err != nil {
		fmt.Println(err)
		return
	}
	txref := factory.Reference()
	userDaoInTx := NewAuthUserDao(&txref)
	id, err = userDaoInTx.Insert(&insertUser)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("insert success!")
	if err = tx.Commit(); err != nil {
		fmt.Println(err)
		return
	}

	_, err = userDao.Delete(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("delete success!")

	// Output:
	// insert success!
	// fetch user from database!
	// abc
	// delete success!
	// insert success!
	// delete success!
}
