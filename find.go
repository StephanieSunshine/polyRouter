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

var table *redis.Client

func  findPeers(geo h3.GeoCoord) string {
	res := 0
	val, err := table.LLen(fmt.Sprintf("%x", h3.FromGeo(geo, res))).Result()
  if err != nil {
    fmt.Println(val, err)
    panic(err)
  }

  //fmt.Println("Count: ", val, "Res: ", res)
  if val < 3 {
    return(findRings(geo))
	}else if val < 31 {
	  return(fmt.Sprintf("%x", h3.FromGeo(geo, res)))
  }else{
	  return(findDrill(geo, val, res))
  }
}

func findRings(geo h3.GeoCoord) string {
  return(fmt.Sprint("To Rings\n", geo))
}

func findDrill(geo h3.GeoCoord, v int64, r int) string {
  val := v
  res := r
  var err error
  for val > 31 {
    if val > 10000 {
      res += 4
    }else if val > 1000 {
      res += 3
    }else if val > 100 {
      res += 2
    }else{
      res += 1
    }
    val, err = table.LLen(fmt.Sprintf("%x", h3.FromGeo(geo, res))).Result()
    if err != nil {
      fmt.Println(val, err)
      panic(err)
    }
  }
  if val < 3 { res -= 1 }
  val, err = table.LLen(fmt.Sprintf("%x", h3.FromGeo(geo, res))).Result()
  if err != nil {
    fmt.Println(val, err)
    panic(err)
  }
  //fmt.Println("Count: ", val, "Res: ", res)
  return(fmt.Sprintf("%x", h3.FromGeo(geo, res)))
}


func main() {
	table = redis.NewClient(&redis.Options{
		Addr:     REDISADR,
    Password: REDISPWD, // no password set
    DB:       0,  // use default DB
	})

	rand.Seed(time.Now().UnixNano())
	for x:=0; x<=100000; x++ {
		geo := h3.GeoCoord{ Latitude: (rand.Float64()*180)-90, Longitude: (rand.Float64()*360)-180 }
    stime := time.Now()
		_ = findPeers(geo)
    fmt.Println(time.Since(stime).Microseconds())
    //fmt.Println(findPeers(geo))
	}
}
