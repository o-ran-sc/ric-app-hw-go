#   Copyright (c) 2019 AT&T Intellectual Property.
#   Copyright (c) 2019 Nokia.
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.

#----------------------------------------------------------
#
#----------------------------------------------------------

FROM nexus3.o-ran-sc.org:10004/o-ran-sc/bldr-ubuntu18-c-go:9-u18.04 as xapp-base
RUN apt-get update -y \
    &&apt-get install -y \
    apt-utils \
    cmake \
    gawk \
    sudo \
    nano \
    jq \
    gettext-base \
    bison \
    flex \
    curl \
    tree

RUN curl -s https://packagecloud.io/install/repositories/o-ran-sc/master/script.deb.sh | bash

# RMR
ARG RMRVERSION=4.2.2
#RUN apt-get install -y rmr=${RMRVERSION} rmr-dev=${RMRVERSION}
RUN wget --content-disposition https://packagecloud.io/o-ran-sc/release/packages/debian/stretch/rmr_${RMRVERSION}_amd64.deb/download.deb && dpkg -i rmr_${RMRVERSION}_amd64.deb
RUN wget --content-disposition https://packagecloud.io/o-ran-sc/release/packages/debian/stretch/rmr-dev_${RMRVERSION}_amd64.deb/download.deb && dpkg -i rmr-dev_${RMRVERSION}_amd64.deb
RUN rm -f rmr_${RMRVERSION}_amd64.deb rmr-dev_${RMRVERSION}_amd64.deb

#
RUN ldconfig

#
ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR "/"
CMD ["/bin/bash"]

#----------------------------------------------------------
#
#----------------------------------------------------------
FROM xapp-base as xapp-base-testbuild


RUN mkdir -p /ws
WORKDIR "/ws"

COPY . /ws
# Module prepare (if go.mod/go.sum updated)
#COPY go.mod /ws
#COPY go.sum /ws
#RUN go mod download

#RUN go build hwApp.go

