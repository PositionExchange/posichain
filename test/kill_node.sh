#!/bin/bash
pkill -9 '^(posichain|soldier|commander|profiler|bootnode)$' | sed 's/^/Killed process: /'
rm -rf db-127.0.0.1-*
