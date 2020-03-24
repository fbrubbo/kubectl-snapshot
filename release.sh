set -e

go test

for GOOS in darwin linux windows; do
  for GOARCH in 386 amd64; do
    output_name="bin/kubectl-snapshot-$GOOS-$GOARCH"
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi
    echo "Building $output_name"
    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name
    if [ $? -ne 0 ]; then
      echo 'An error has occurred! Aborting the script execution...'
      exit 1
    fi
  done
done

chmod +x ./sh/*
cp ./sh/* ./bin/


# TO TEST LOCALLY
# 1. ./install.sh
# 2. git add -A && git commit -m "Message: #issue" && git push
# 4. create new release in github and download the Source tar.gz
# 5. get the sha256 hash using https://emn178.github.io/online-tools/sha256_checksum.html
# 6. replace all 'uri' and 'sha256' entries in resource-snapshot.yaml
# 7. kubectl krew uninstall resource-snapshot
# 8. kubectl krew install --manifest=resource-snapshot.yaml
# 9. kubectl resource-snapshot -v
# 10. kubectl krew uninstall resource-snapshot
# 10. git add -A && git commit -m "v0.1.3" && git push (just to have the new uri and sha256 saved)


# TO DEPLOY ON KREW
# 1. submite a pull request to https://github.com/kubernetes-sigs/krew-index using resource-snapshot.yaml content
# 2. kubectl krew uninstall resource-snapshot
# 3. kubectl krew update
# 4. kubectl krew install resource-snapshot


