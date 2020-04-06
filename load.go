package main

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"math/rand"
	"time"
	"github.com/uber/h3-go"
)

const REDISADR = "localhost:6379"
const REDISPWD = ""

type Location struct {
	Lat float64
	Lon float64
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     REDISADR,
		Password: REDISPWD, // no password set
		DB:       15,  // use default DB
	})
	
	table := redis.NewClient(&redis.Options{
		Addr:     REDISADR,
		Password: REDISPWD, // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
	rand.Seed(time.Now().UnixNano())

	for c:=0; c<100000; c ++ {
	address := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255)+1, rand.Intn(256), rand.Intn(256), rand.Intn(255)+1)
	loc := Location{(rand.Float64()*180)-90, (rand.Float64()*360)-180}
	
	err = client.HMSet(address, []string{"lat", fmt.Sprintf("%f", loc.Lat), "lon", fmt.Sprintf("%f", loc.Lon)}).Err()

	fmt.Println(address, loc)
	geo := h3.GeoCoord{ Latitude: loc.Lat, Longitude: loc.Lon }

	for r:=0; r<16; r++ {
		err := table.LPush(fmt.Sprintf("%x", h3.FromGeo(geo, r)), address).Err()
		if err != nil {
			panic(err)
		}
	}


	}
}
