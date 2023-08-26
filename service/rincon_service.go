package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"san_francisco/config"
	"san_francisco/model"
	"san_francisco/utils"
	"strconv"
	"strings"
	"time"
)

var rinconRetries = 0
var rinconHost = "http://rincon" + ":" + config.RinconPort

func RegisterRincon() {
	var portInt, _ = strconv.Atoi(config.Port)
	rinconBody, _ := json.Marshal(map[string]interface{}{
		"name":         "SanFrancisco",
		"version":      config.Version,
		"url":          "http://sanfrancisco:" + config.Port,
		"port":         portInt,
		"status_email": config.StatusEmail,
	})
	// Azure Container App deployment
	var ContainerAppEnvDNSSuffix = os.Getenv("CONTAINER_APP_ENV_DNS_SUFFIX")
	if ContainerAppEnvDNSSuffix != "" {
		utils.SugarLogger.Infoln("Found Azure Container App environment variables, using internal DNS suffix: " + ContainerAppEnvDNSSuffix)
		rinconHost = "http://rincon.internal." + ContainerAppEnvDNSSuffix
		rinconBody, _ = json.Marshal(map[string]interface{}{
			"name":         "SanFrancisco",
			"version":      config.Version,
			"url":          "http://sanfrancisco." + ContainerAppEnvDNSSuffix,
			"port":         portInt,
			"status_email": config.StatusEmail,
		})
	}

	responseBody := bytes.NewBuffer(rinconBody)
	res, err := http.Post(rinconHost+"/services", "application/json", responseBody)
	if err != nil {
		utils.SugarLogger.Errorln(err.Error())
		if rinconRetries < 15 {
			rinconRetries++
			if rinconRetries%2 == 0 {
				rinconHost = "http://localhost" + ":" + config.RinconPort
				utils.SugarLogger.Errorln("failed to register with rincon, retrying with \"http://localhost\" in 5s...")
			} else {
				rinconHost = "http://rincon" + ":" + config.RinconPort
				utils.SugarLogger.Errorln("failed to register with rincon, retrying with \"http://rincon\" in 5s...")
			}
			time.Sleep(time.Second * 5)
			RegisterRincon()
		} else {
			utils.SugarLogger.Fatalln("failed to register with rincon after 15 attempts, terminating program...")
		}
	} else {
		defer res.Body.Close()
		if res.StatusCode == 200 {
			json.NewDecoder(res.Body).Decode(&config.Service)
			println("===========================================")
			println("service info: " + config.Service.Name + " v" + config.Service.Version)
			println("service info: " + config.Service.URL)
			println()
		} else {
			utils.SugarLogger.Errorln("Failed to register with Rincon! Status code: " + strconv.Itoa(res.StatusCode))
			// print body
			buf := new(bytes.Buffer)
			buf.ReadFrom(res.Body)
			newStr := buf.String()
			println("===========================================")
			println(newStr)
			println()
		}
		utils.SugarLogger.Infoln("Registered service with Rincon! Service ID: " + strconv.Itoa(config.Service.ID))
		RegisterRinconRoute("/sanfrancisco")
		GetRinconServiceInfo()
	}
}

func RegisterRinconRoute(route string) {
	rinconBody, _ := json.Marshal(map[string]string{
		"route":        route,
		"service_name": "SanFrancisco",
	})
	responseBody := bytes.NewBuffer(rinconBody)
	_, err := http.Post(rinconHost+"/routes", "application/json", responseBody)
	if err != nil {
		utils.SugarLogger.Errorln(err.Error())
	}
	utils.SugarLogger.Infoln("Registered route " + route)
}

func GetRinconServiceInfo() {
	var service model.Service
	rinconClient := http.Client{}
	req, _ := http.NewRequest("GET", rinconHost+"/routes/match/rincon", nil)
	res, err := rinconClient.Do(req)
	if err != nil {
		utils.SugarLogger.Errorln(err.Error())
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		json.NewDecoder(res.Body).Decode(&service)
	}
	println("===========================================")
	println("Rincon service info: " + service.Name + " v" + service.Version)
	println("Rincon service info: " + service.URL)
	println()
	config.RinconService = service
}

func MatchRoute(traceparent string, route string, requestID string) model.Service {
	var service model.Service
	queryRoute := strings.ReplaceAll(route, "/", "<->")
	rinconClient := http.Client{}
	req, _ := http.NewRequest("GET", rinconHost+"/routes/match/"+queryRoute, nil)
	req.Header.Set("Request-ID", requestID)
	req.Header.Set("traceparent", traceparent)
	req.Header.Add("Content-Type", "application/json")
	res, err := rinconClient.Do(req)
	if err != nil {
		utils.SugarLogger.Errorln(err.Error())
		return service
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		json.NewDecoder(res.Body).Decode(&service)
	}
	return service
}
