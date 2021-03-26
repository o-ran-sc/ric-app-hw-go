.. This work is licensed under a Creative Commons Attribution 4.0 International License.
.. SPDX-License-Identifier: CC-BY-4.0
.. Copyright (c) 2021 Samsung Electronics Co., Ltd. All Rights Reserved.Copyright (C) 2021

============================================================================================ 
HW-go xAPP (golang)
============================================================================================ 
-------------------------------------------------------------------------------------------- 
User's Guide 
-------------------------------------------------------------------------------------------- 
 
Introduction 
============================================================================================ 

The RIC platform provides set of functions as part of xAPP golang Framework that the xAPPs can use to accomplish their tasks.
This xAPP is envisioned to provide python xAPP developers, examples of implementing these sets of functions.
Note, HW-go xAPP does not address/implement any RIC Usecases. 

HW-go xAPP Features 
============================================================================================ 

RIC Platform provides many Frameworks, APIs and libraries to aid the development of xAPPs. All xAPPs will have some custom
processing functional logic core to the xApp and some additional non-functional platform related processing using 
these APIs and libraries. This xAPP attempts to show the usage of such additional platform processing using xapp RIC framework APIs and libraries.


The HW-go xAPP demonstrates how a golang based xApp uses the A1, and E2 interfaces and persistent database read-write operations.

