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

package main

import (
	"github.com/edgexfoundry/device-rest-rfrain"
	"github.com/edgexfoundry/device-rest-rfrain/driver"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
)

const (
	serviceName string = "edgex-device-rest-rfrain"
)

func main() {
	sd := driver.RestDriver{}
	startup.Bootstrap(serviceName, device_rest_rfrain.Version, &sd)
}
