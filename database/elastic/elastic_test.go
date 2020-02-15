package elastic

import (
	"math/rand"
	"time"
)

const indexName = "item_increase_log"

type ItemIncreaseLog struct {
	UserID   uint64 `elastic:""`
	Increase int64
	Recent   int64
	Time     int64
	ItemID   uint32
	Type     uint32
}

func NewPointerIncreaseLogs(length int) (logs []*ItemIncreaseLog) {

	currency := rand.Int63n(1000) + 1000

	for i := 0; i < length; i++ {
		time.Sleep(time.Microsecond * time.Duration(rand.Intn(5)+1))
		increase := 100 - rand.Int63n(200)
		logs = append(logs, &ItemIncreaseLog{
			UserID:   userID,
			Increase: increase,
			Recent:   currency,
			Time:     time.Now().UnixNano() / int64(time.Millisecond),
			ItemID:   itemID,
			Type:     uint32(rand.Intn(10)),
		})
		currency += increase
	}

	return logs
}

func NewStructIncreaseLogs(length int) (logs []ItemIncreaseLog) {

	currency := rand.Int63n(1000) + 1000

	for i := 0; i < length; i++ {
		time.Sleep(time.Microsecond * time.Duration(rand.Intn(5)+1))
		increase := 100 - rand.Int63n(200)
		logs = append(logs, ItemIncreaseLog{
			UserID:   userID,
			Increase: increase,
			Recent:   currency,
			Time:     time.Now().UnixNano() / int64(time.Millisecond),
			ItemID:   itemID,
			Type:     uint32(rand.Intn(10)),
		})
		currency += increase
	}

	return
}

const (
	userID = 10000000
	itemID = 1
)

var client *Elastic

func init() {

	rand.Seed(time.Now().UnixNano())

	var err error
	client, err = New(Option{Address: "http://192.168.1.6:9200"})
	if err != nil {
		panic(err)
	}
}
