# !/bin/sh
echo "== .piebot -s stop =="
./piebot -s stop
echo "== go generate =="
go generate
echo "== go build =="
go build
echo "== .piebot =="
./piebot

