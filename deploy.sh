#!/bin/bash -e

# Deploys onto swami-direct (the HA-IP doesn't forward SSH traffic)
DEPLOYHOST=swami-direct.blender.cloud
DEPLOYPATH=/opt/  # end in slash!
SSH="ssh -o ClearAllForwardings=yes"

echo "======== Building a statically-linked svnman"
bash docker/build-via-docker.sh linux
source ./docker/_version.sh

echo "======== Deploying onto $DEPLOYHOST"
rsync -e "$SSH" -va docker/$PREFIX-linux.tar.gz $DEPLOYHOST:
$SSH $DEPLOYHOST -t <<EOT
set -ex
cd $DEPLOYPATH
tar zxvf ~/$PREFIX-linux.tar.gz
rm -f svn-manager
ln -s $PREFIX svn-manager
EOT

echo
echo "======== Restarting service"
$SSH $DEPLOYHOST -t "sudo systemctl restart svn-manager.service"

echo "======== Deploy done."
