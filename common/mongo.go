package common
import (
	"gopkg.in/mgo.v2"
	"container/list"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"fmt"
)


//err = collection.Update(bson.M{"name": "ccc"}, bson.M{"$set": bson.M{"name": "ddd"}})
//_, err = collection.RemoveAll(bson.M{"name": "Ale”})
//err = collection.Find(bson.M{"phone": "456"}).One(&result)

func NewMongo(config *Config) *mgo.Database {
	//[mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
	session, err := mgo.Dial(config.MongDB_Addr)
	if err != nil {
		panic(err)
	}

	return session.DB(config.MongDB_DB)
}

func GetAllStatusApiOfUser(db *mgo.Database, uid string) *list.List {
	c := db.C("status_api")

	result := list.New()
	api := StatusAPI{}

	it := c.Find(bson.M{"userId": uid, "monitor": true}).Iter()
	for it.Next(&api) {
		var newAPI = fill(&api)

		switch api.Type{
		case "GET":
			b, _ := json.Marshal(api.Content)
			json.Unmarshal(b, &newAPI.GetContent)
			//fmt.Printf("%s, %+v\n", har.Name, newHar.GetContent)
		case "POST":
			b, _ := json.Marshal(api.Content)
			json.Unmarshal(b, &newAPI.PostContent)
			//fmt.Printf("%s, %+v\n", har.Name, newHar.PostContent)
		default:
		}
		//fmt.Println("-------")

		result.PushBack(newAPI)
	}

	return result
}


func GetAllStatusApi(db *mgo.Database) *list.List {
	c := db.C("status_api")

	result := list.New()
	api := StatusAPI{}

	it := c.Find(bson.M{"monitor": true}).Iter()
	for it.Next(&api) {
		var newAPI = fill(&api)

		switch api.Type{
		case "GET":
			b, _ := json.Marshal(api.Content)
			json.Unmarshal(b, &newAPI.GetContent)
		case "POST":
			b, _ := json.Marshal(api.Content)
			json.Unmarshal(b, &newAPI.PostContent)
		default:
		}

		result.PushBack(newAPI)
	}

	return result
}


func GetStatusApi(db *mgo.Database, apiId string) *StatusAPI {
	c := db.C("status_api")

	api := StatusAPI{}

	it := c.Find(bson.M{ "id":apiId}).Iter()
	for it.Next(&api) {
		switch api.Type{
		case "GET":
			b, _ := json.Marshal(api.Content)
			json.Unmarshal(b, &api.GetContent)
		case "POST":
			b, _ := json.Marshal(api.Content)
			json.Unmarshal(b, &api.PostContent)
		default:
		}
	}

	return &api
}

func SaveStatusAPILog(db *mgo.Database, log *StatusAPILog){
	c := db.C("status_api_log")
	err:=c.Insert(log)
	if err!=nil{
		fmt.Printf("error when save log: %s", err)
	}
}

func fill(api *StatusAPI) *StatusAPI {
	var newAPI = &StatusAPI{}
	newAPI.Id = api.Id
	newAPI.Monitor = api.Monitor
	newAPI.Cron = api.Cron
	newAPI.Name = api.Name
	newAPI.UserId = api.UserId
	newAPI.Date = api.Date
	newAPI.Type=api.Type

	return newAPI
}