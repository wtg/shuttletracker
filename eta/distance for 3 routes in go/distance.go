package main

import "fmt"
import "os"
import "io/ioutil"
import "encoding/json"
import "math"

func main(){
	routeFile,err :=os.Open("route.json")
	fmt.Println(err)
	var route []map[string]interface{}
	routeValue,_ := ioutil.ReadAll(routeFile)
	json.Unmarshal(routeValue,&route)
	// fmt.Println(len(route))
	//fmt.Println(route[0]["points"].([]interface{})[0].(map[string]interface{})["latitude"])
	// var coords [3][][]float64 
	var distances [3]float64 
	for i:=0; i<len(route);i++{
	 	fmt.Println(route[i]["id"],route[i]["name"])

	 	rtmp:=route[i]["points"].([]interface{})
	 	// coords[i]=make([][]float64,len(rtmp))
	 	var dis float64=0
	 	var lat,lon float64
	 	for j:=0; j<len(rtmp);j++{
	 		if j==0{
		 		lat=rtmp[j].(map[string]interface{})["latitude"].(float64)/360*2*math.Pi
		 		lon=rtmp[j].(map[string]interface{})["longitude"].(float64)/360*2*math.Pi
		 	}else{
		 		lat1:=rtmp[j].(map[string]interface{})["latitude"].(float64)/360*2*math.Pi
		 		lon1:=rtmp[j].(map[string]interface{})["longitude"].(float64)/360*2*math.Pi
		 		dlon:=lon1-lon
		 		dlat:=lat1-lat
		 		a:=math.Pow(math.Sin(dlat/2),2) + math.Cos(lat1)*math.Cos(lat)*math.Pow(math.Sin(dlon/2),2)
		 		c := 2*math.Asin(math.Sqrt(a))
		 		dis+=c*6371/1.6
		 		lat=lat1
		 		lon=lon1
		 		

		 	}

	 		//fmt.Println(lat,lon)
	 	}
	 	fmt.Println(dis)
	 	distances[i]=dis
	}
	fmt.Println()
	fmt.Println(distances)
	
}