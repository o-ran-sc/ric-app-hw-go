.. This work is licensed under a Creative Commons Attribution 4.0 International License.
.. SPDX-License-Identifier: CC-BY-4.0
.. Copyright (c) 2021 Samsung Electronics Co., Ltd. All Rights Reserved.Copyright (C) 2021

.. _Deployment Guide:

Onboarding and Deployment of Hw-go xApp
=======================================

.. contents::
  :depth: 3
  :local:

Onboarding of hw-go using dms_cli tool
--------------------------------------

`dms_cli` offers rich set of command line utility to onboard *hw-go* xapp
to `chartmuseme`.

First checkout the `hw-go <https://gerrit.o-ran-sc.org/r/admin/repos/ric-app/hw-go>`_ repository from gerrit.

.. code-block:: bash

  $ git clone "https://gerrit.o-ran-sc.org/r/ric-app/hw-go"


`hw-go` has following folder structure

.. code-block:: bash

  ├── Dockerfile
  ├── INFO.yaml
  ├── LICENSES.txt
  ├── README.md
  ├── config
  │   ├── config-file.json           // descriptor for hw-go
  │   ├── schema.json                // schema for controls section of descriptor
  │   └── uta_rtg.rt                 // local route file
  ├── docs
  ├── go.mod
  ├── go.sum
  └── hwApp.go


For onboarding `hw-go` make sure that `dms_cli` and helm3 is installed. One can follow `documentation <https://docs.o-ran-sc.org/projects/o-ran-sc-it-dep/en/latest/installation-guides.html#ric-applications>`_ to
configure `dms_cli`.

Once `dms_cli` is availabe we can proceed to onboarding procedure.

configure the `export CHART_REPO_URL` to point `chartmuseme`.

.. code-block:: bash

  $export CHART_REPO_URL=http://<service-ricplt-xapp-onboarder-http.ricplt>:8080

check if `dms_cli` working fine.

.. code-block:: bash

  $ dms_cli health
  True

Now move to `config` folder to initiate onboarding.

.. code-block:: bash

  $ cd config
  $ dms_cli onboard --config_file_path=config-file.json --shcema_file_path=schema.json
  httpGet:
    path: '{{ index .Values "readinessProbe" "httpGet" "path" | toJson }}'
    port: '{{ index .Values "readinessProbe" "httpGet" "port" | toJson }}'
  initialDelaySeconds: '{{ index .Values "readinessProbe" "initialDelaySeconds" | toJson }}'
  periodSeconds: '{{ index .Values "readinessProbe" "periodSeconds" | toJson }}'

  httpGet:
    path: '{{ index .Values "livenessProbe" "httpGet" "path" | toJson }}'
    port: '{{ index .Values "livenessProbe" "httpGet" "port" | toJson }}'
  initialDelaySeconds: '{{ index .Values "livenessProbe" "initialDelaySeconds" | toJson }}'
  periodSeconds: '{{ index .Values "livenessProbe" "periodSeconds" | toJson }}'

  {
      "status": "Created"
  }

Check if `hw-go` is onborded

.. code-block:: bash

  $ curl --location --request GET "http://<appmgr>:32080/onboard/api/v1/charts"  --header 'Content-Type: application/json'

  {
      "hw-go": [
          {
              "name": "hw-go",
              "version": "1.0.0",
              "description": "Standard xApp Helm Chart",
              "apiVersion": "v1",
              "appVersion": "1.0",
              "urls": [
                  "charts/hw-go-1.0.0.tgz"
              ],
              "created": "2021-06-24T18:57:41.98056196Z",
              "digest": "14a484d9a394ed34eab66e5241ec33e73be8fa70a2107579d19d037f2adf57a0"
          }
      ]
  }

If we would wish to download the charts then we can perform following curl operation :

.. code-block:: bash

  curl --location --request GET "http://<appmgr>:32080/onboard/api/v1/charts/xapp/hw-go/ver/1.0.0"  --header 'Content-Type: application/json' --output hw-go.tgz

Now the onboarding is done.

Deployment of hw-go
-------------------

Once charts are available we can deploy the the `hw-go` using following curl command :

.. code-block:: bash

  $ curl --location --request POST "http://<appmgr>:32080/appmgr/ric/v1/xapps"  --header 'Content-Type: application/json'  --data-raw '{"xappName": "hw-go", "helmVersion": "1.0.0"}'
  {"instances":null,"name":"hw-go","status":"deployed","version":"1.0"}

Deployment will be done in `ricxapp` ns :

.. code-block:: bash

  # kubectl get po -n ricxapp
  NAME                             READY   STATUS    RESTARTS   AGE
  ricxapp-hw-go-55ff7549df-kpj6k   1/1     Running   0          2m

  # kubectl get svc -n ricxapp
  NAME                         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)             AGE
  aux-entry                    ClusterIP   IP1             <none>        80/TCP,443/TCP      73d
  service-ricxapp-hw-go-http   ClusterIP   IP2             <none>        8080/TCP            103m
  service-ricxapp-hw-go-rmr    ClusterIP   IP3             <none>        4560/TCP,4561/TCP   103m

Now we can query to appmgr to get list of all the deployed xapps :

.. code-block:: bash

  # curl http://service-ricplt-appmgr-http.ricplt:8080/ric/v1/xapps | jq .
    % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                  Dload  Upload   Total   Spent    Left  Speed
  100   347  100   347    0     0    578      0 --:--:-- --:--:-- --:--:--   579
  [
    {
      "instances": [
        {
          "ip": "service-ricxapp-hw-go-rmr.ricxapp",
          "name": "hw-go-55ff7549df-kpj6k",
          "policies": [
            1
          ],
          "port": 4560,
          "rxMessages": [
            "RIC_SUB_RESP",
            "A1_POLICY_REQ",
            "RIC_HEALTH_CHECK_REQ"
          ],
          "status": "running",
          "txMessages": [
            "RIC_SUB_REQ",
            "A1_POLICY_RESP",
            "A1_POLICY_QUERY",
            "RIC_HEALTH_CHECK_RESP"
          ]
        }
      ],
      "name": "hw-go",
      "status": "deployed",
      "version": "1.0"
    }
  ]

Logs from `hw-go` :

.. code-block:: bash

  # kubectl  logs ricxapp-hw-go-55ff7549df-kpj6k -n ricxapp
  {"ts":1624562552123,"crit":"INFO","id":"hw-app","mdc":{"time":"2021-06-24T19:22:32"},"msg":"Using config file: config/config-file.json"}
  {"ts":1624562552124,"crit":"INFO","id":"hw-app","mdc":{"CONTAINER_NAME":"","HOST_NAME":"","PID":"6","POD_NAME":"","SERVICE_NAME":"","SYSTEM_NAME":"","time":"2021-06-24T19:22:32"},"msg":"Serving metrics on: url=/ric/v1/metrics namespace=ricxapp"}
  {"ts":1624562552133,"crit":"INFO","id":"hw-app","mdc":{"CONTAINER_NAME":"","HOST_NAME":"","PID":"6","POD_NAME":"","SERVICE_NAME":"","SYSTEM_NAME":"","time":"2021-06-24T19:22:32"},"msg":"Register new counter with opts: {ricxapp SDL Stored The total number of stored SDL transactions map[]}"}
  {"ts":1624562552133,"crit":"INFO","id":"hw-app","mdc":{"CONTAINER_NAME":"","HOST_NAME":"","PID":"6","POD_NAME":"","SERVICE_NAME":"","SYSTEM_NAME":"","time":"2021-06-24T19:22:32"},"msg":"Register new counter with opts: {ricxapp SDL StoreError The total number of SDL store errors map[]}"}
  1624562552 6/RMR [INFO] ric message routing library on SI95 p=0 mv=3 flg=00 (fd4477a 4.5.2 built: Jan 21 2021)
  {"ts":1624562552140,"crit":"INFO","id":"hw-app","mdc":{"CONTAINER_NAME":"","HOST_NAME":"","HWApp":"0.0.1","PID":"6","POD_NAME":"","SERVICE_NAME":"","SYSTEM_NAME":"","time":"2021-06-24T19:22:32"},"msg":"new rmrClient with parameters: ProtPort=0 MaxSize=0 ThreadType=0 StatDesc=RMR LowLatency=false FastAck=false Policies=[]"}
  {"ts":1624562552140,"crit":"INFO","id":"hw-app","mdc":{"CONTAINER_NAME":"","HOST_NAME":"","HWApp":"0.0.1","PID":"6","POD_NAME":"","SERVICE_NAME":"","SYSTEM_NAME":"","time":"2021-06-24T19:22:32"},"msg":"Register new counter with opts: {ricxapp RMR Transmitted The total number of transmited RMR messages map[]}"}
  {"ts":1624562552140,"crit":"INFO","id":"hw-app","mdc":{"CONTAINER_NAME":"","HOST_NAME":"","HWApp":"0.0.1","PID":"6","POD_NAME":"","SERVICE_NAME":"","SYSTEM_NAME":"","time":"2021-06-24T19:22:32"},"msg":"Register new counter with opts: {ricxapp RMR Received The total number of received RMR messages map[]}"}
  {"ts":1624562552140,"crit":"INFO","id":"hw-app","mdc":{"CONTAINER_NAME":"","HOST_NAME":"","HWApp":"0.0.1","PID":"6","POD_NAME":"","SERVICE_NAME":"","SYSTEM_NAME":"","time":"2021-06-24T19:22:32"},"msg":"Register new counter with opts: {ricxapp RMR TransmitError The total number of RMR transmission errors map[]}"}
  {"ts":1624562552140,"crit":"INFO","id":"hw-app","mdc":{"CONTAINER_NAME":"","HOST_NAME":"","HWApp":"0.0.1","PID":"6","POD_NAME":"","SERVICE_NAME":"","SYSTEM_NAME":"","time":"2021-06-24T19:22:32"},"msg":"Register new counter with opts: {ricxapp RMR ReceiveError The total number of RMR receive errors map[]}"}
  {"ts":1624562552140,"crit":"INFO","id":"hw-app","mdc":{"CONTAINER_NAME":"","HOST_NAME":"","HWApp":"0.0.1","PID":"6","POD_NAME":"","SERVICE_NAME":"","SYSTEM_NAME":"","time":"2021-06-24T19:22:32"},"msg":"Xapp started, listening on: :8080"}
  {"ts":1624562552140,"crit":"INFO","id":"hw-app","mdc":{"CONTAINER_NAME":"","HOST_NAME":"","HWApp":"0.0.1","PID":"6","POD_NAME":"","SERVICE_NAME":"","SYSTEM_NAME":"","time":"2021-06-24T19:22:32"},"msg":"rmrClient: Waiting for RMR to be ready ..."}
  {"ts":1624562553140,"crit":"INFO","id":"hw-app","mdc":{"CONTAINER_NAME":"","HOST_NAME":"","HWApp":"0.0.1","PID":"6","POD_NAME":"","SERVICE_NAME":"","SYSTEM_NAME":"","time":"2021-06-24T19:22:33"},"msg":"rmrClient: RMR is ready after 1 seconds waiting..."}
  {"ts":1624562553141,"crit":"INFO","id":"hw-app","mdc":{"CONTAINER_NAME":"","HOST_NAME":"","HWApp":"0.0.1","PID":"6","POD_NAME":"","SERVICE_NAME":"","SYSTEM_NAME":"","time":"2021-06-24T19:22:33"},"msg":"xApp ready call back received"}
  1624562553 6/RMR [INFO] sends: ts=1624562553 src=service-ricxapp-hw-go-rmr.ricxapp:0 target=localhost:4591 open=0 succ=0 fail=0 (hard=0 soft=0)
  1624562553 6/RMR [INFO] sends: ts=1624562553 src=service-ricxapp-hw-go-rmr.ricxapp:0 target=localhost:4560 open=0 succ=0 fail=0 (hard=0 soft=0)
  1624562553 6/RMR [INFO] sends: ts=1624562553 src=service-ricxapp-hw-go-rmr.ricxapp:0 target=service-ricplt-a1mediator-rmr.ricplt:4562 open=0 succ=0 fail=0 (hard=0 soft=0)
  RMR is ready now ...
  {"ts":1624562557140,"crit":"INFO","id":"hw-app","mdc":{"CONTAINER_NAME":"","HOST_NAME":"","HWApp":"0.0.1","PID":"6","POD_NAME":"","SERVICE_NAME":"","SYSTEM_NAME":"","time":"2021-06-24T19:22:37"},"msg":"Application='hw-go' is not ready yet, waiting ..."}
  {"ts":1624562562141,"crit":"INFO","id":"hw-app","mdc":{"CONTAINER_NAME":"","HOST_NAME":"","HWApp":"0.0.1","PID":"6","POD_NAME":"","SERVICE_NAME":"","SYSTEM_NAME":"","time":"2021-06-24T19:22:42"},"msg":"Application='hw-go' is not ready yet, waiting ..."}

Here we are done with the onboaring and deployment of `hw-go`.

