#!/bin/sh

echo $PRIVATE_KEY > private_key;
base64 -d private_key > id_rsa;
echo "" >> id_rsa;
mkdir ~/.ssh/;
cp id_rsa ~/.ssh/id_rsa;
eval `ssh-agent`;
chmod 600 ~/.ssh/id_rsa;
ssh-add;
echo "StrictHostKeyChecking no" >> ~/.ssh/config;
git config --global user.email $USER_EMAIL;
git clone git@github.com:$ORGANIZATION/$REPOSITORY.git;
git --git-dir=/$REPOSITORY/.git merge $SRC_BRANCH;
git --git-dir=/$REPOSITORY/.git push origin $DEST_BRANCH;
echo ""
echo "Merged successfully"
echo "<-- END -->{\"success\": true}"