package common
import (
	"gopkg.in/mgo.v2"
	"container/list"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
)


//err = collection.Update(bson.M{"name": "ccc"}, bson.M{"$set": bson.M{"name": "ddd"}})
//_, err = collection.RemoveAll(bson.M{"name": "Ale”})
//err = collection.Find(bson.M{"phone": "456"}).One(&result)

func NewMongo() *mgo.Database {
	//[mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
	session, err := mgo.Dial("192.168.100.186:20301")
	if err != nil {
		panic(err)
	}

	return session.DB("moklr")
}

func GetAllStatusApiOfUser(db *mgo.Database, uid string) *list.List {
	c := db.C("status_api")

	result := list.New()
	har := Har{}

	it := c.Find(bson.M{"userId": uid, "monitor": true}).Iter()
	for it.Next(&har) {
		var newHar = fill(&har)

		switch har.Type{
		case "GET":
			b, _ := json.Marshal(har.Content)
			json.Unmarshal(b, &newHar.GetContent)
			//fmt.Printf("%s, %+v\n", har.Name, newHar.GetContent)
		case "POST":
			b, _ := json.Marshal(har.Content)
			json.Unmarshal(b, &newHar.PostContent)
			//fmt.Printf("%s, %+v\n", har.Name, newHar.PostContent)
		default:
		}
		//fmt.Println("-------")

		result.PushBack(newHar)
	}

	return result
}


func GetAllStatusApi(db *mgo.Database) *list.List {
	c := db.C("status_api")

	result := list.New()
	har := Har{}

	it := c.Find(bson.M{"monitor": true}).Iter()
	for it.Next(&har) {
		var newHar = fill(&har)

		switch har.Type{
		case "GET":
			b, _ := json.Marshal(har.Content)
			json.Unmarshal(b, &newHar.GetContent)
		case "POST":
			b, _ := json.Marshal(har.Content)
			json.Unmarshal(b, &newHar.PostContent)
		default:
		}

		result.PushBack(newHar)
	}

	return result
}

func fill(har *Har) *Har{
	var newHar = &Har{}
	newHar.HarId = har.HarId
	newHar.Name = har.Name
	newHar.UserId = har.UserId
	newHar.CollectionId = har.CollectionId
	newHar.Date = har.Date
	newHar.Type=har.Type

	return newHar
}