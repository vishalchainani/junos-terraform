#!/bin/bash

set -e

COMMON_GROUP_MAX_DEPTH="${COMMON_GROUP_MAX_DEPTH:-4}"
COMMON_GROUP_MAX_COUNT="${COMMON_GROUP_MAX_COUNT:-12}"
COMMON_GROUP_MIN_BENEFIT="${COMMON_GROUP_MIN_BENEFIT:-2}"
COMMON_GROUP_MIN_NEW_PATHS="${COMMON_GROUP_MIN_NEW_PATHS:-2}"
RESET_PAYLOAD_CACHE="${RESET_PAYLOAD_CACHE:-1}"

if [[ "$RESET_PAYLOAD_CACHE" == "1" ]]; then
	rm -f ansible_files/.jtaf_host_payloads.json
fi

jtaf-xml2yaml -x ../evpn-vxlan-dc/dc1/dc1-*leaf* ../evpn-vxlan-dc/dc1/dc1-*spine* -j ansible-provider-junos-vqfx-ansible-role/trimmed_schema.json -d ansible_files --common-host-groups-max-depth "$COMMON_GROUP_MAX_DEPTH" --common-host-groups-max-count "$COMMON_GROUP_MAX_COUNT" --common-host-groups-min-benefit "$COMMON_GROUP_MIN_BENEFIT" --common-host-groups-min-new-paths "$COMMON_GROUP_MIN_NEW_PATHS"

jtaf-xml2yaml -x ../evpn-vxlan-dc/dc1/dc1-*firewall* -j ansible-provider-junos-srx-ansible-role/trimmed_schema.json -d ansible_files --common-host-groups-max-depth "$COMMON_GROUP_MAX_DEPTH" --common-host-groups-max-count "$COMMON_GROUP_MAX_COUNT" --common-host-groups-min-benefit "$COMMON_GROUP_MIN_BENEFIT" --common-host-groups-min-new-paths "$COMMON_GROUP_MIN_NEW_PATHS"
