#!/bin/bash

[ -z $1 ] && echo "Usage: $0 tarfile" && exit 1

[ -d bin ] && echo "You MUST delete directory bin manually" && exit 1

tar -xf $1
