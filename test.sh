#!/bin/bash

cpu=$(cat /proc/cpuinfo | grep processor | wc -l)
cpu=`expr $((cpu)) - 1`
echo $cpu