package try

import (
	"errors"
	"fmt"
	"testing"
)

func TestTry(t *testing.T) {
	type Err1 struct {
		error
	}
	type Err2 struct {
		error
	}

	//Try1
	Try(func() {
		fmt.Println("Try1 start")
		panic(Err1{error: errors.New("error1")})
	}).Catch(Err1{}, func(err error) {
		fmt.Println("Try1 Err1 Catch:", err.Error())
	}).Catch(Err2{}, func(err error) {
		fmt.Println("Try1 Err2 Catch:", err.Error())
	}).Finally(func() {
		fmt.Println("Try1 done")
	})

	//Try2
	Try(func() {
		fmt.Println("Try2 start")
		panic(Err2{error: errors.New("error2")})
	}).Catch(Err1{}, func(err error) {
		fmt.Println("Try2 Err1 Catch:", err.Error())
	}).CatchAll(func(err error) {
		fmt.Println("Try2 Catch All:", err.Error())
	}).Finally(func() {
		fmt.Println("Try2 done")
	})

	//Try3
	Try(func() {
		fmt.Println("Try3 start")
	}).Catch(Err1{}, func(err error) {
		fmt.Println("Try3 Err1 Catch:", err.Error())
	}).Finally(func() {
		fmt.Println("Try3 done")
	})
}
