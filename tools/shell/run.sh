#!/bin/bash
cd gateway && $(dirname $(pwd))/bin/gateway > gateway.log 2>&1 &
cd game && $(dirname $(pwd))/bin/game > game.log 2>&1 &
cd door && $(dirname $(pwd))/bin/door > door.log 2>&1 &