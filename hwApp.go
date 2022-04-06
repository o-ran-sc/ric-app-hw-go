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

==================================================================================
*/
package main

import (
	"encoding/json"

	"gerrit.o-ran-sc.org/r/ric-plt/alarm-go.git/alarm"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/clientmodel"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
)

type HWApp struct {
	stats map[string]xapp.Counter
}

var (
	A1_POLICY_QUERY      = 20013
	POLICY_QUERY_PAYLOAD = "{\"policy_type_id\":20000}"
	reqId                = int64(1)
	seqId                = int64(1)
	funId                = int64(1)
	actionId             = int64(1)
	actionType           = "report"
	subsequestActioType  = "continue"
	timeToWait           = "w10ms"
	direction            = int64(0)
	procedureCode        = int64(27)
	xappEventInstanceID  = int64(1234)
	typeOfMessage        = int64(1)
	subscriptionId       = ""
	hPort                = int64(8080)
	rPort                = int64(4560)
	clientEndpoint       = clientmodel.SubscriptionParamsClientEndpoint{Host: "service-ricxapp-hw-go-rmr.ricxapp", HTTPPort: &hPort, RMRPort: &rPort}
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

func (e *HWApp) getEnbList() ([]*xapp.RNIBNbIdentity, error) {
	enbs, err := xapp.Rnib.GetListEnbIds()

	if err != nil {
		xapp.Logger.Error("err: %s", err)
		return nil, err
	}

	xapp.Logger.Info("List for connected eNBs :")
	for index, enb := range enbs {
		xapp.Logger.Info("%d. enbid: %s", index+1, enb.InventoryName)
	}
	return enbs, nil
}

func (e *HWApp) getGnbList() ([]*xapp.RNIBNbIdentity, error) {
	gnbs, err := xapp.Rnib.GetListGnbIds()

	if err != nil {
		xapp.Logger.Error("err: %s", err)
		return nil, err
	}

	xapp.Logger.Info("List of connected gNBs :")
	for index, gnb := range gnbs {
		xapp.Logger.Info("%d. gnbid : %s", index+1, gnb.InventoryName)
	}
	return gnbs, nil
}

func (e *HWApp) getnbList() []*xapp.RNIBNbIdentity {
	nbs := []*xapp.RNIBNbIdentity{}

	if enbs, err := e.getEnbList(); err == nil {
		nbs = append(nbs, enbs...)
	}

	if gnbs, err := e.getGnbList(); err == nil {
		nbs = append(nbs, gnbs...)
	}
	return nbs
}

func (e *HWApp) sendSubscription(meid string) {

	xapp.Logger.Info("sending subscription request for meid : %s", meid)

	subscriptionParams := clientmodel.SubscriptionParams{
		ClientEndpoint: &clientEndpoint,
		Meid:           &meid,
		RANFunctionID:  &funId,
		SubscriptionDetails: clientmodel.SubscriptionDetailsList([]*clientmodel.SubscriptionDetail{
			{
				ActionToBeSetupList: clientmodel.ActionsToBeSetup{
					&clientmodel.ActionToBeSetup{
						ActionDefinition: clientmodel.ActionDefinition([]int64{1, 2, 3, 4}),
						ActionID:         &actionId,
						ActionType:       &actionType,
						SubsequentAction: &clientmodel.SubsequentAction{
							SubsequentActionType: &subsequestActioType,
							TimeToWait:           &timeToWait,
						},
					},
				},
				EventTriggers:       clientmodel.EventTriggerDefinition([]int64{1, 2, 3, 4}),
				XappEventInstanceID: &xappEventInstanceID,
			},
		}),
	}

	b, err := json.MarshalIndent(subscriptionParams, "", "  ")

	if err != nil {
		xapp.Logger.Error("Json marshaling failed : %s", err)
		return
	}

	xapp.Logger.Info("*****body: %s ", string(b))

	resp, err := xapp.Subscription.Subscribe(&subscriptionParams)

	if err != nil {
		xapp.Logger.Error("subscription failed (%s) with error: %s", meid, err)

		// subscription failed, raise alarm
		err := xapp.Alarm.Raise(8086, alarm.SeverityCritical, meid, "subscriptionFailed")
		if err != nil {
			xapp.Logger.Error("Raising alarm failed with error %v", err)
		}

		return
	}
	xapp.Logger.Info("Successfully subcription done (%s), subscription id : %s", meid, *resp.SubscriptionID)
}

func (e *HWApp) xAppStartCB(d interface{}) {
	xapp.Logger.Info("xApp ready call back received")

	// get the list of all NBs
	nbList := e.getnbList()

	// send subscription request to each of the NBs
	for _, nb := range nbList {
		e.sendSubscription(nb.InventoryName)
	}
}

func (e *HWApp) handleRICIndication(ranName string, r *xapp.RMRParams) {
	// update metrics for indication message
	e.stats["RICIndicationRx"].Inc()
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

	// RIC INDICATION message
	case "RIC_INDICATION":
		xapp.Logger.Info("Received RIC Indication message")
		e.handleRICIndication(msg.Meid.RanName, msg)

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
	// Defind metrics counter that the xapp provides
	metrics := []xapp.CounterOpts{
		{
			Name: "RICIndicationRx",
			Help: "Total number of RIC Indication message received",
		},
	}

	hw := HWApp{
		stats: xapp.Metric.RegisterCounterGroup(metrics, "hw_go"), // register counter
	}
	hw.Run()
}
