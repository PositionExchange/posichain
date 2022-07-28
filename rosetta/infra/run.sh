#!/usr/bin/env bash
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
DATA="$DIR/data"
LOGS="$DATA/logs"
DATA_NAME="${DATA_NAME:=posichain_sharddb_0}"

MAINNET_22816573_SNAPSHOT="release:pub.posichain.org/mainnet.min.22816573/posichain_sharddb_0"

case "$NETWORK" in
mainnet)
  CONFIG_PATH="-c /root/posichain-mainnet.conf"
  ;;
mainnet-22816573)
  CONFIG_PATH="-c /root/posichain-mainnet.conf"
  rclone -P -L sync $MAINNET_22816573_SNAPSHOT $DATA/$DATA_NAME --transfers=64
  ;;
testnet)
  # TODO add testnet conf
  CONFIG_PATH="-c /root/posichain-testnet.conf"
  ;;
*)
  echo "unknown network"
  exit 1
  ;;
esac

if [ "$MODE" = "offline" ]; then
  BASE_ARGS=(--datadir "$DATA" --log.dir "$LOGS" --run.offline)
else
  BASE_ARGS=(--datadir "$DATA" --log.dir "$LOGS")
fi

mkdir -p "$LOGS"
echo -e NODE ARGS: \" $CONFIG_PATH "$@" "${BASE_ARGS[@]}" \"
echo "NODE VERSION: $($DIR/posichain --version)"

"$DIR/posichain" $CONFIG_PATH "$@" "${BASE_ARGS[@]}"
