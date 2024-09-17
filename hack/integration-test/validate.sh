#!/bin/bash

set -ex

echo "Waiting for configuration ISO URL to be set on ImageClusterInstall ibi-test/ibi-test"

for i in {1..120}; do
  url=$(oc get -n ibi-test dataimage ostest-extraworker -ojsonpath='{.spec.url}' || true)
  if [[ -n "$url" ]]; then
    break
  fi
  sleep 1
done

if [[ -z "$url" ]]; then
  echo "ERROR: configurationImageURL on ImageClusterInstall ibi-test/ibi-test was not set within 60 seconds"
  exit 1
fi

echo "Renaming url to match the new image-based-install-operator route as we can't reach service directly"
new_url=$(echo "$url" | sed 's|image-based-install-config.image-based-install-operator.svc:8000|images-image-based-install-operator.apps.dev.redhat.com|')

curl --insecure --output "test.iso" "$new_url"
file "test.iso" | grep -q "ISO 9660 CD-ROM"
