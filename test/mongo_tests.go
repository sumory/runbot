package main

import (
	"fmt"
	"github.com/sumory/runbot/common"
)

func main() {


	db := common.NewMongo()
	hars := common.GetAllStatusAPIOfUser(db, "db46161c-917a-40d8-b1fd-242e7cc8f4b3")
	fmt.Println(hars.Len())

	for e := hars.Front(); e != nil; e = e.Next() {
		element := e.Value
		h := element.(*common.Har)

		if h.Type=="POST"{

			fmt.Printf("%+v\n", h.PostContent.PostData.Params)
		}

//		e,_,r := common.RunHar(h)
//		fmt.Println(e,r)
		//fmt.Println("---------------------")
	}
}