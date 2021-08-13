#!/bin/bash

cmd=$1
ip=$2
port=$3
user=$4
password=$5

# construct connection string
c_ip=$ip
[ ! $c_ip ] && c_ip="127.0.0.1"

c_port=$port
[ ! $c_port ] && c_port=27017

[ $user ] && c_user="$user:$password@"

cnnstr="mongodb://$c_user$c_ip:$c_port/${PROJ_NAME}b?authSource=admin"

# construct mongodump/mongorestore args
[ $ip ]       && ip="-h $ip"
[ $port ]     && port="--port $port"
[ $user ]     && user="-u $user"
[ $password ] && password="--password=$password"

# get to work
case "$cmd" in
    make)
        mongo "$cnnstr" --quiet --eval '
            db.order.aggregate([
                {$match: {
                    status: "ok",
                }},
                {$group: {
                    _id: "$userid",
                    bill: {$sum: "$amount"},
                }},
                {$sort: {
                    bill: -1,
                }},
                {$lookup: {
                    from: "refund",
                    localField: "_id",
                    foreignField: "_id",
                    as: "refund",
                }},
                {$unwind: "$refund"},
                {$project: {
                    _id:  "$refund.code",
                    bill: 1,
                }},
                {$out: "refund_take"},
            ]);
        '

        mongodump $ip $port $user $password --authenticationDatabase=admin -d ${PROJ_NAME}b -c refund_take --gzip --archive=refund_take.db
        ;;

    write)
        read -p "Are you sure to write refund_take?[Yes/No]" x
        [ "$x" != "Yes" ] && exit 1

        mongorestore $ip $port $user $password --authenticationDatabase=admin --drop --gzip --archive=refund_take.db
        ;;

    *)
        echo "$0 {make|write} [ip] [port] [user] [password]"
esac

