#!/bin/bash

set -e

if [ ! -e /run/secrets/kubernetes.io/serviceaccount ] && [ ! -e /dev/kmsg ]; then
    echo "ERROR: Require --privileged flag to run Seal"
    exit 1
fi

# make it work with docker desktop on M1
# https://github.com/rancher/rancher/issues/35930
if [ -f /sys/fs/cgroup/cgroup.controllers ]; then
  echo "[$(date -Iseconds)] [CgroupV2 Fix] Evacuating Root Cgroup ..."
	# move the processes from the root group to the /init group,
  # otherwise writing subtree_control fails with EBUSY.
  mkdir -p /sys/fs/cgroup/init
  xargs -rn1 < /sys/fs/cgroup/cgroup.procs > /sys/fs/cgroup/init/cgroup.procs || :
  # enable controllers
  sed -e 's/ / +/g' -e 's/^/+/' <"/sys/fs/cgroup/cgroup.controllers" >"/sys/fs/cgroup/cgroup.subtree_control"
  echo "[$(date -Iseconds)] [CgroupV2 Fix] Done"
fi

# start
exec "$@"
