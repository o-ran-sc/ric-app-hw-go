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

var (
	A1_POLICY_QUERY      = 20013
	POLICY_QUERY_PAYLOAD = "{\"policy_type_id\":20000}"
)

func (e *HWApp) sendPolicyQuery() {
	xapp.Logger.Info("Invoked method to send  policy query message")

	// prepare and send policy query message over RMR
	rmrParams := new(xapp.RMRParams)
	rmrParams.Mtype = A1_POLICY_QUERY // A1_POLICY_QUERY
	rmrParams.Payload = []byte(POLICY_QUERY_PAYLOAD)

	// send rmr message
	flg := xapp.Rmr.SendMsg(rmrParams)

	if flg {
		xapp.Logger.Info("Successfully sent policy query message over RMR")
	} else {
		xapp.Logger.Info("Failed to send policy query message over RMR")
	}
}

func (e *HWApp) ConfigChangeHandler(f string) {
	xapp.Logger.Info("Config file changed")
}

func (e *HWApp) xAppStartCB(d interface{}) {
	xapp.Logger.Info("xApp ready call back received")
}

func (e *HWApp) Consume(msg *xapp.RMRParams) (err error) {
	id := xapp.Rmr.GetRicMessageName(msg.Mtype)

	xapp.Logger.Info("Message received: name=%s meid=%s subId=%d txid=%s len=%d", id, msg.Meid.RanName, msg.SubId, msg.Xid, msg.PayloadLen)

	switch id {
	// policy request handler
	case "A1_POLICY_REQUEST":
		xapp.Logger.Info("Recived policy instance list")

	// health check request
	case "RIC_HEALTH_CHECK_REQ":
		xapp.Logger.Info("Received health check request")

	default:
		xapp.Logger.Info("Unknown message type '%d', discarding", msg.Mtype)
	}

	defer func() {
		xapp.Rmr.Free(msg.Mbuf)
		msg.Mbuf = nil
	}()
	return
}

func (e *HWApp) Run() {

	// set MDC
	xapp.Logger.SetMdc("HWApp", "0.0.1")

	// set config change listener
	xapp.AddConfigChangeListener(e.ConfigChangeHandler)

	// register callback after xapp ready
	xapp.SetReadyCB(e.xAppStartCB, true)

	// reading configuration from config file
	waitForSdl := xapp.Config.GetBool("db.waitForSdl")

	// start xapp
	xapp.RunWithParams(e, waitForSdl)

}

func main() {
	hw := HWApp{}
	hw.Run()
}
