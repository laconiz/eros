package main

import (
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"log"
)

func main() {

	plan, err := watch.Parse(map[string]interface{}{
		"type":   "keyprefix",
		"prefix": "oceanus/",
	})
	if err != nil {
		panic(err)
	}

	plan.Handler = func(u uint64, i interface{}) {
		if pairs, ok := i.(api.KVPairs); ok {
			if len(pairs) > 0 {
				for _, pair := range pairs {
					log.Printf("%#v", pair)
				}
			} else {
				log.Println("empty pairs")
			}
		} else {
			log.Println("nil pairs")
		}
	}

	plan.Run("127.0.0.1:8500")
}
