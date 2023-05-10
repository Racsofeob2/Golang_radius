package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
	//"regexp"	

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	
)

func main() {
	filename := "PyCharm/EXPORT.csv"
	// Load the CSV file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Parse the CSV data
	reader := csv.NewReader(file)
	reader.Comma = ';'
	ssidMap := make(map[string]string)
	for {
		row, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatal(err)
		}
		mac := strings.TrimSpace(row[0])
		ssid := strings.TrimSpace(row[1])
		ssidMap[mac] = ssid
	}

	handler := func(w radius.ResponseWriter, r *radius.Request) {
	
vlanID := "99" //identificador VLAN IPs puntos Wi-Fi

    // Allow access for the general AP
    if strings.HasPrefix(r.RemoteAddr.String(), "192.168." + vlanID + ".") {

        // Get the client's MAC address from the RADIUS packet
        macAddr := rfc2865.CallingStationID_GetString(r.Packet)
	log.Printf(macAddr)

        // Check if the client's MAC is allowed in the MAC-SSID table
        allowedSSID := ssidMap[macAddr]
	ssid := strings.Split(rfc2865.CalledStationID_GetString(r.Packet), ":")[1]

        if ssid != allowedSSID {

            code := radius.CodeAccessReject
            log.Printf("Writing %v to %v", code, r.RemoteAddr)
            w.Write(r.Response(code))
            return

        }else{


        // Client's MAC is allowed
        code := radius.CodeAccessAccept
        log.Printf("Writing %v to %v", code, r.RemoteAddr)
        w.Write(r.Response(code))
        return

    	}
}
    // Access not allowed for non-AP requests
    code := radius.CodeAccessReject
    log.Printf("Writing %v to %v", code, r.RemoteAddr)
    w.Write(r.Response(code))
}

	server := radius.PacketServer{
		Handler:      radius.HandlerFunc(handler),
		SecretSource: radius.StaticSecretSource([]byte(`secret`)),
		 	
	}

	log.Printf("Starting server on :1812")




	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}