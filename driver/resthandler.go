//
// Copyright (c) 2019 Intel Corporation
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
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"encoding/json"

	"github.com/edgexfoundry/device-sdk-go/pkg/models"
	sdk "github.com/edgexfoundry/device-sdk-go/pkg/service"
	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
)

const (
	apiResourceRoute  = clients.ApiBase + "/alert"
	handlerContextKey = "RFRainHandler"
)

type RestHandler struct {
	service     *sdk.Service
	logger      logger.LoggingClient
	asyncValues chan<- *models.AsyncValues
	rfRain      *RFRainClient
}

func NewRestHandler(service *sdk.Service, logger logger.LoggingClient, asyncValues chan<- *models.AsyncValues, rf *RFRainClient) *RestHandler {
	handler := RestHandler{
		service:     service,
		logger:      logger,
		asyncValues: asyncValues,
		rfRain:      rf,
	}

	return &handler
}

func (handler RestHandler) Start() error {
	if err := handler.service.AddRoute(apiResourceRoute, handler.addContext(deviceHandler), http.MethodPost); err != nil {
		return fmt.Errorf("unable to add required route: %s: %s", apiResourceRoute, err.Error())
	}

	handler.logger.Info(fmt.Sprintf("Route %s added.", apiResourceRoute))

	return nil
}

func (handler RestHandler) addContext(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	// Add the context with the handler so the endpoint handling code can get back to this handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), handlerContextKey, handler)
		next(w, r.WithContext(ctx))
	})
}

func (handler RestHandler) processAsyncRequest(writer http.ResponseWriter, request *http.Request) {
	var alerts Alerts

	handler.logger.Debug(fmt.Sprintf("Alerts being received"))
	rawResp, err := handler.readBodyAsString(writer, request)
	if err != nil {
		handler.logger.Error(fmt.Sprintf("Incoming alert ignored. Unable to read request body: %s", err.Error()))
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal([]byte(rawResp), &alerts)
	if err != nil {
		handler.logger.Error(fmt.Sprintf("unable to unmarshal RFRain alerts: %s", err))
		return
	} else {
		for _, alertEvent := range alerts.Alert {
			handler.logger.Debug(fmt.Sprintf("Incoming reader alert: %+v", alertEvent))
			_, err := handler.service.GetDeviceByName(alertEvent.Readerid)
			if err != nil {
				handler.logger.Error(fmt.Sprintf("Incoming alert ignored. Device '%s' not found", alertEvent.Readerid))
				http.Error(writer, fmt.Sprintf("Device '%s' not found", alertEvent.Readerid), http.StatusNotFound)
				return
			}
			handler.logger.Info(fmt.Sprintf("Device '%s' found", alertEvent.Readerid))

			asyncValues := &models.AsyncValues{
				DeviceName:    alertEvent.Readerid,
				CommandValues: handler.createReadings(alertEvent),
			}
			handler.asyncValues <- asyncValues
		}
	}
	return
}

func (handler RestHandler) createReadings(alertEvent AlertEvent) []*models.CommandValue {
	// #TODO make this more dynamic based on device profile
	var cvs = make([]*models.CommandValue, 11)
	var result = &models.CommandValue{}
	var timestamp = time.Now().UnixNano()
	result = models.NewStringValue("tagnumb",timestamp,alertEvent.Tagnumb)
	cvs[0] = result
	result = models.NewStringValue("tagname",timestamp,alertEvent.Tagname)
	cvs[1] = result
	result = models.NewStringValue("subzone",timestamp,alertEvent.Subzone)
	cvs[2] = result
	result = models.NewStringValue("SS",timestamp,alertEvent.Ss)
	cvs[3] = result
	result = models.NewStringValue("readername",timestamp,alertEvent.Readername)
	cvs[4] = result
	result = models.NewStringValue("groupname",timestamp,alertEvent.Groupname)
	// cvs[5] = result
	// result = models.NewStringValue("readerid",timestamp,alertEvent.Readerid)
	// cvs[6] = result
	// result = models.NewStringValue("data",timestamp,alertEvent.Data)
	cvs[5] = result
	result = models.NewStringValue("zone",timestamp,alertEvent.Zone)
	cvs[6] = result
	result = models.NewStringValue("location",timestamp,alertEvent.Location)
	cvs[7] = result
	result = models.NewStringValue("current_status",timestamp,alertEvent.Current_Status)
	cvs[8] = result
	result = models.NewStringValue("current_access_utc",timestamp,alertEvent.Current_Access_Utc)
	cvs[9] = result
	result = models.NewStringValue("current_access",timestamp,alertEvent.Current_Access)
	cvs[10] = result
	handler.logger.Debug(fmt.Sprintf("Command values: %v", cvs))
	return cvs
}

func (handler RestHandler) readBodyAsString(writer http.ResponseWriter, request *http.Request) (string, error) {
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return "", err
	}

	if len(body) == 0 {
		return "", fmt.Errorf("no request body provided")
	}

	return string(body), nil
}

func deviceHandler(writer http.ResponseWriter, request *http.Request) {
	handler, ok := request.Context().Value(handlerContextKey).(RestHandler)
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Bad context pass to handler"))
		return
	}

	handler.processAsyncRequest(writer, request)
}
