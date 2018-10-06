# !/bin/sh

postResult=$(curl "http://127.0.0.1:8085/wallet" -d "{\"symbol\":\"$1\",\"txid\":\"$2\"}")
if [[ "$postResult" != "accept success" ]];then
    minuteRange=5
    filenamepart1=$(date "+%Y%m%d%H")
    minute=$(date "+%M")
    filenamepart2=$[(10#$minute)/minuteRange*minuteRange]
    if [[ $filenamepart2 -lt 10 ]];then
        filenamepart2="0$filenamepart2"
    fi
    echo $filenamepart2
    filename=$filenamepart1$filenamepart2.$minuteRange
    filepath=/Users/dumh/walletnotifydemo/txs/$1
    [[ ! -d $filepath ]] && mkdir -p $filepath
    echo "$1,$2" >> $filepath/$filename
fi

