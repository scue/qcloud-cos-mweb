#!/bin/bash
source /anaconda3/bin/activate /anaconda3/envs/py27
# edit ~/.cos.conf
# OR, coscmd config -a <secret_id> -s <secret_key> -b <bucket> -r <region>
./qcloud-cos-upload "$@"
