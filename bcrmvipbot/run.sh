# !/bin/sh
echo "== .bcrmvipbot -s stop =="
./bcrmvipbot -s stop
echo "== go generate =="
go generate
echo "== go build =="
go build
echo "== .bcrmvipbot =="
./bcrmvipbot

