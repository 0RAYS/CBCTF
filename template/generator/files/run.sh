#!/bin/bash

# $1: team_id
# $2: base64(flag)
flag=`echo $2 | base64 -d`

# 生成附件
python3 generator.py $1 $flag