.. This work is licensed under a Creative Commons Attribution 4.0 International License.
.. SPDX-License-Identifier: CC-BY-4.0
.. Copyright (c) 2021 Samsung Electronics Co., Ltd. All Rights Reserved.


Installation Guide
==================

.. contents::
   :depth: 3
   :local:

Abstract
--------

This document describes how to install the reference HW-go xAPP.


Introduction
------------

This document provides guidelines on how to install and configure the HW-go xAPP in various environments/operating modes.
The audience of this document is assumed to have good knowledge in RIC Platform.


Preface
-------
This xAPP can be run directly as a Linux binary, as a docker image, or in a pod in a Kubernetes environment.  The first
two can be used for dev testing. The last option is how an xAPP is deployed in the RAN Intelligent Controller environment.
This document covers all three methods.




Software Installation and Deployment
------------------------------------
The build process assumes a Linux environment with go >= 1.15  and  has been tested on Ubuntu. For building docker images,
the Docker environment must be present in the system.


Build Process
~~~~~~~~~~~~~
The HW-go xAPP can be either tested as a Linux binary or as a docker image.
   1. **Linux binary**:
      TBD

   2. **Docker Image**: From the root of the repository, run   *docker --no-cache build -t <image-name> ./* .


Deployment
~~~~~~~~~~

End to end deployment of `hw-go` can be referred at :

  :ref:`Deployment Guide`.

Testing
--------

Unit tests TBD
