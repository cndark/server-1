#!/bin/bash

echo "edit first!"
exit 1

df|grep /dev/vdb
[ $? -eq 0 ] && echo "already mounted" && exit 1

fdisk -l|grep /dev/vdb
[ $? -ne 0 ] && echo "no data-disk to mount" && exit 1

mkfs.ext4 /dev/vdb
[ ! -d /data ] && mkdir /data
mount /dev/vdb /data
chown game:game /data
rm -rf /data/*

cat /etc/fstab|grep /dev/vdb
[ ! $? -eq 0 ] && echo '/dev/vdb /data ext4 defaults 0 0' >> /etc/fstab

