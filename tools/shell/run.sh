#!/bin/bash
cd gateway && $(dirname $(pwd))/bin/gateway > log.txt 2>&1 &
cd game && $(dirname $(pwd))/bin/game > log.txt 2>&1 &
cd door && $(dirname $(pwd))/bin/door > log.txt 2>&1 &