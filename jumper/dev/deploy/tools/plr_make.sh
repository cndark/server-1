#!/bin/bash

echo "edit first!" && exit 1

addr="http://192.168.0.101:12001/gm/service"
plrid="u1-1000000"

curl -d "plrid=$plrid&key=w.plr&xp=155833524" $addr
curl -d "plrid=$plrid&key=w.res&res_k=10100&res_v=99999999&res_k=10203&res_v=99999999&res_k=10305&res_v=9999999&res_k=10302&res_v=500000" $addr
curl -d "plrid=$plrid&key=w.res&res_k=30102&res_v=999&res_k=30103&res_v=999" $addr
curl -d "plrid=$plrid&key=w.res&res_k=30210&res_v=999&res_k=30211&res_v=999&res_k=30212&res_v=999" $addr
