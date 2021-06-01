/*
==================================================================================
  Copyright (c) 2020 Samsung

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

   This source code is part of the near-RT RIC (RAN Intelligent Controller)
   platform project (RICP).
==================================================================================
*/
package main

import (
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
)

type HWApp struct {
}

func (e *HWApp) ConfigChangeHandler(f string) {
	xapp.Logger.Info("Config file changed")
}


func (e *HWApp) xAppStartCB(d interface{}) {
	xapp.Logger.Info("xApp ready call back received")
}

func (e *HWApp) Consume(rp *xapp.RMRParams) (err error) {
	return
}

func (e *HWApp) Run() {

	// set MDC
	xapp.Logger.SetMdc("HWApp", "0.0.1")

	// set config change listener
	xapp.AddConfigChangeListener(e.ConfigChangeHandler)

	// register callback after xapp ready
	xapp.SetReadyCB(e.xAppStartCB, true)

	xapp.RunWithParams(e, false)

}

func main() {
	hw := HWApp{}
	hw.Run()
}
