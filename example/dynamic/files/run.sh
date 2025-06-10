#!/bin/bash

# $1: team_id
# $2: base64(base64(flag1),base64(flag2),...)

# 生成附件
python3 generator.py $1 $2