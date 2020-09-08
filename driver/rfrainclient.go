//
// Copyright (c) 2020 IOTech Systems
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package driver

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/edgexfoundry/device-sdk-go/pkg/service"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"io/ioutil"
	"net/http"
)

type SessionKeyResp struct {
	Request string
	Success bool
	Results struct {
		Sessionkey string
		Userlevel  string
	}
	Message string
}

type ResultResponse struct {
	Category string
	Request  string
	Success  bool
	Message  string
	Results  [] struct {
		Groupid   string
		Groupname string
		Email     []string
		Api       []string
	}
}

type WorkedResponse struct {
	Category string
	Request  string
	Success  bool
	Message  string
	Results  struct {
		Worked bool
		Msg    string
	}
}

type Alerts struct {
	Alert []AlertEvent
}

type AlertEvent struct {
	Tagnumb string
	Tagname string
	//Detectstat	string
	Subzone            string
	Ss                 string
	Readername         string
	Groupname          string
	Readerid           string
	// Data               string
	Zone               string
	Location           string
	Current_Status     string
	Current_Access_Utc string
	Current_Access     string
}

type RFRainClient struct {
	// TODO  - organize config into structs
	SessionKey                 string
	User                       string
	Password                   string
	Company                    string
	SessionKeyURL              string
	StartMonitorURL            string
	InvalidateURL              string
	StartAlertEngineURL        string
	StopAlertEngineURL         string
	StartAlertURL              string
	StopAlertURL               string
	DefineAlertGroupURL        string
	ListAlertGroupURL          string
	DefineAlertURL             string
	AddAlertAPIEntryToGroupURL string
	AlertRecEndpointURL        string
	GroupName                  string
	AlertId                    int
	Logger                     logger.LoggingClient
}

func (c *RFRainClient) loadCredentials() error {
	c.Logger.Debug("Loading RFRain credentials")
	c.User = sdk.DriverConfigs()["User"]
	c.Password = sdk.DriverConfigs()["Password"]
	c.Company = sdk.DriverConfigs()["Company"]
	c.SessionKeyURL = sdk.DriverConfigs()["SessionKeyURL"]
	c.StartMonitorURL = sdk.DriverConfigs()["StartMonitoringURL"]
	c.InvalidateURL = sdk.DriverConfigs()["InvalidateURL"]
	c.StartAlertEngineURL = sdk.DriverConfigs()["StartAlertEngineURL"]
	c.StopAlertEngineURL = sdk.DriverConfigs()["StopAlertEngineURL"]
	c.StartAlertURL = sdk.DriverConfigs()["StartAlertURL"]
	c.StopAlertURL = sdk.DriverConfigs()["StopAlertURL"]
	c.DefineAlertGroupURL = sdk.DriverConfigs()["DefineAlertGroupURL"]
	c.ListAlertGroupURL = sdk.DriverConfigs()["ListAlertGroupURL"]
	c.DefineAlertURL = sdk.DriverConfigs()["DefineAlertURL"]
	c.DefineAlertGroupURL = sdk.DriverConfigs()["DefineAlertGroupURL"]
	c.AddAlertAPIEntryToGroupURL = sdk.DriverConfigs()["AddAlertAPIEntryToGroupURL"]
	c.AlertRecEndpointURL = sdk.DriverConfigs()["AlertRecEndpointURL"]
	c.GroupName = sdk.DriverConfigs()["GroupName"]
	c.Logger.Debug("RFRain credentials Loaded")
	return nil
}

func NewRFRainClient(lc logger.LoggingClient) *RFRainClient {
	client := new(RFRainClient)
	client.Logger = lc
	client.loadCredentials()
	return client
}

func (c *RFRainClient) StartSession() {
	sessionOk := c.getSessionKey()
	if sessionOk {
		jsonData := map[string]string{"sessionkey": c.SessionKey}
		jsonValue, _ := json.Marshal(jsonData)
		fmt.Sprintf("%s", jsonValue)
		// if the alert does not exist
		if !c.listAlertGroups() {
			c.Logger.Debug("alert group not defined... adding alert group: %s", c.GroupName)
			// #TODO
			//define alert group
			//define alert
			//add alert api entry to group
		}
		//start alert engine
		c.executeSessionCommand(c.StartAlertEngineURL)
		// start alerts
		c.executeSessionCommandResultResponse(c.StartAlertURL)
	}
}

func (c *RFRainClient) getSessionKey() bool {
	jsonData := map[string]string{"email": base64.StdEncoding.EncodeToString([]byte(c.User)), "cname": c.Company, "password": base64.StdEncoding.EncodeToString([]byte(c.Password))}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post(c.SessionKeyURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		c.Logger.Error(fmt.Sprintf("Problem getting RFRain session key: %s", err))
		return false
	} else {
		sessKeyResp := SessionKeyResp{}
		rawResp, _ := ioutil.ReadAll(response.Body)
		err = json.Unmarshal(rawResp, &sessKeyResp)
		if err != nil {
			c.Logger.Error(fmt.Sprintf("Unable to unmarshal RFRain session key response; session key may have expired: %s", err))
			return false
		} else {
			if sessKeyResp.Success {
				c.SessionKey = sessKeyResp.Results.Sessionkey
				c.Logger.Info(fmt.Sprintf("RFRain session key successfully obtained"))
				c.Logger.Debug(fmt.Sprintf("RFRain session key:  %s", c.SessionKey))
			} else {
				c.Logger.Error(fmt.Sprintf("Unsuccessful login attempt: %s", sessKeyResp.Message))
				return false
			}
		}
	}
	return true
}

func (c *RFRainClient) listAlertGroups() bool {
	listAlertGroup, err := c.executeSessionCommandResultResponse(c.ListAlertGroupURL)
	if err != nil {
		c.Logger.Error(fmt.Sprintf("Unable to get existing alert groups. %s", err))
		// #TODO where to use panics
		panic(err)
	} else {
		c.Logger.Debug(fmt.Sprintf("Looking for existing alert group:  %s\n", c.GroupName))
		for _, g := range listAlertGroup.Results {
			c.Logger.Debug(fmt.Sprintf("Found alert group: %s", g.Groupname))
			if g.Groupname == c.GroupName {
				c.Logger.Debug(fmt.Sprintf("Existing group name located.  No need to add alert group %s\n", g.Groupname))
				return true
			}
		}
	}
	return false
}

func (c *RFRainClient) EndSession() error {
	// #TODO - handle case of when session id goes invalid
	// stop alerts
	_, err := c.executeSessionCommandResultResponse(c.StopAlertURL)
	if err != nil {
		c.Logger.Error(fmt.Sprintf("Problem stopping alerts: %s", err))
		return err
	}

	// stop alert engine
	err = c.executeSessionCommand(c.StopAlertEngineURL)
	if err != nil {
		c.Logger.Error(fmt.Sprintf("Problem stopping alert engine: %s", err))
		return err
	}

	// invalidate the session
	_, err = c.executeSessionCommandResultResponse(c.InvalidateURL)
	if err != nil {
		c.Logger.Error(fmt.Sprintf("Problem invalidating the session: %s", err))
		return err
	}
	c.Logger.Info(fmt.Sprintf("Successful session end and key invalidation"))
	return nil
}

func (c *RFRainClient) executeSessionCommandResultResponse(commandURL string) (*ResultResponse, error) {
	response := &ResultResponse{}
	jsonData := map[string]string{"sessionkey": c.SessionKey}
	jsonValue, _ := json.Marshal(jsonData)
	resp, err := http.Post(commandURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		c.Logger.Error(fmt.Sprintf("Problem executing command URL: %s.  Error:  %s", commandURL, err))
		defer resp.Body.Close()
		return response, err
	}
	rawResp, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(rawResp, &response)
	if err != nil {
		c.Logger.Error(fmt.Sprintf("Problem marshalling response on command URL: %s.  Error: %s", commandURL, err))
		defer resp.Body.Close()
		return response, err
	}
	if !response.Success {
		c.Logger.Error(fmt.Sprintf("Unsuccessful call on command URL: %s.  Error %s", commandURL, response.Message))
		defer resp.Body.Close()
		return response, errors.New(response.Message)
	}
	c.Logger.Debug(fmt.Sprintf("Successful command execution: %s", commandURL))
	defer resp.Body.Close()
	return response, nil
}

func (c *RFRainClient) executeSessionCommand(commandURL string) error {
	response := &WorkedResponse{}
	jsonData := map[string]string{"sessionkey": c.SessionKey}
	jsonValue, _ := json.Marshal(jsonData)
	resp, err := http.Post(commandURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		c.Logger.Error(fmt.Sprintf("Problem executing command URL: %s.  Error:  %s", commandURL, err))
		defer resp.Body.Close()
		return err
	}
	rawResp, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(rawResp, &response)
	if err != nil {
		c.Logger.Error(fmt.Sprintf("Problem marshalling response on command URL: %s.  Error: %s", commandURL, err))
		defer resp.Body.Close()
		return err
	}
	if !response.Success {
		c.Logger.Error(fmt.Sprintf("Unsuccessful call on command URL: %s.  Error %s", commandURL, response.Message))
		defer resp.Body.Close()
		return errors.New(response.Message)
	}
	c.Logger.Debug(fmt.Sprintf("Successful command execution: %s", commandURL))
	defer resp.Body.Close()
	return nil
}
