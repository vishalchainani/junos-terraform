# Release Note: `--grouping-hosts-file` Knob for `jtaf-xml2yaml`

## Overview

The `--grouping-hosts-file` flag is a new required argument added to the `jtaf-xml2yaml` command. It gives operators explicit control over how Ansible inventory groups and `group_vars` directories are structured when converting Junos XML configs into Ansible host/group variable files.

---

## Why It Was Added

In real deployments, devices are logically grouped by role (e.g., `spine`, `leaf`, `firewall`). With `--grouping-hosts-file`, the operator fully owns the group topology and the tool produces a predictable, stable inventory structure across repeated runs.

---

## The Grouping Hosts File Format

The file uses the same INI syntax as an Ansible inventory file. Two types of sections are supported:

| Section syntax | Purpose |
|---|---|
| `[all]` | Lists every host that should appear in the `[all]` section of the generated inventory |
| `[groupname]` | Lists the hosts that belong to a flat inventory group called `groupname` |

Lines beginning with `#` are treated as comments. Each host or group token is one entry per line.

---

## Simple Flat Example — Single Data Center (DC1)

One grouping hosts file is created per data center. Suppose DC1 contains seven QFX devices. Create `dc1.grouping.hosts`:

```ini
[all]
dc1-borderleaf1
dc1-borderleaf2
dc1-leaf1
dc1-leaf2
dc1-leaf3
dc1-spine1
dc1-spine2

[borderleaf]
dc1-borderleaf1
dc1-borderleaf2

[leaf]
dc1-leaf1
dc1-leaf2
dc1-leaf3

[spine]
dc1-spine1
dc1-spine2
```

Then run:

```bash
jtaf-xml2yaml \
  -j ansible-provider-junos-vqfx-evpn-vxlan/trimmed_schema.json \
  -x examples/evpn-vxlan-dc/dc1/*{spine,leaf}*.xml \
  -d ansible-evpn-vxlan-deploy \
  --grouping-hosts-file ansible-evpn-vxlan-deploy/dc1.grouping.hosts
```

The tool writes:

```
ansible-evpn-vxlan-deploy/
├── inventory.ini
├── group_vars/
│   ├── all.yaml              ← config keys shared by ALL seven DC1 devices
│   ├── borderleaf/
│   │   └── all.yaml          ← keys shared only by dc1-borderleaf1 and dc1-borderleaf2
│   ├── leaf/
│   │   └── all.yaml          ← keys shared only by the three leaf devices
│   └── spine/
│       └── all.yaml          ← keys shared only by the two DC1 spine devices
└── host_vars/
    ├── dc1-borderleaf1.yaml  ← unique delta for this device
    ├── dc1-leaf1.yaml
    └── ...
```

The generated `inventory.ini` will contain the `[borderleaf]`, `[leaf]`, and `[spine]` group sections exactly as declared in the grouping file.

---

## Multiple Data Centers in One Inventory Directory

A key design point is that two separate `jtaf-xml2yaml` runs — each with its own `--grouping-hosts-file` — can target the **same `-d` output directory**. The tool merges inventory sections without clobbering existing groups.

Continuing the example above, DC2 contains only spine devices. Create `dc2.grouping.hosts`:

```ini
[all]
dc2-spine1
dc2-spine2

[spine]
dc2-spine1
dc2-spine2
```

Run the second conversion targeting the same directory:

```bash
jtaf-xml2yaml \
  -j ansible-provider-junos-vqfx-evpn-vxlan/trimmed_schema.json \
  -x examples/evpn-vxlan-dc/dc2/*spine*.xml \
  -d ansible-evpn-vxlan-deploy \
  --grouping-hosts-file ansible-evpn-vxlan-deploy/dc2.grouping.hosts
```

After both runs, `ansible-evpn-vxlan-deploy/inventory.ini` contains all groups from both data centers:

```ini
[all]
dc1-borderleaf1
dc1-borderleaf2
dc1-leaf1
dc1-leaf2
dc1-leaf3
dc1-spine1
dc1-spine2
dc2-spine1
dc2-spine2

[borderleaf]
dc1-borderleaf1
dc1-borderleaf2

[leaf]
dc1-leaf1
dc1-leaf2
dc1-leaf3

[spine]
dc1-spine1
dc1-spine2
dc2-spine1
dc2-spine2
```

And `group_vars/spine/all.yaml` contains config shared across all four spine devices from both data centers — with any spine-specific differences between DC1 and DC2 falling through to the individual `host_vars` files.

---

## Variable Deduplication Behaviour

`jtaf-xml2yaml` performs a three-level deduplication pass driven entirely by the group structure declared in the grouping hosts file:

1. **`group_vars/all.yaml`** — config keys whose value is identical across every host currently tracked in the output directory.
2. **`group_vars/<group>/all.yaml`** — keys that are identical within a group but differ from the global shared value.
3. **`host_vars/<hostname>.yaml`** — only the residual, host-unique differences.

If a host's config conflicts with its group siblings, the conflicting keys are placed directly in that host's `host_vars` file rather than in the group, preserving correctness.

---

## Validation and Error Handling

The parser validates the grouping file at startup and exits with a clear message for:
- Missing or unreadable file
- Malformed `[section]` headers
- Duplicate section names
- Hosts in a group that were not provided via `-x` (when `--strict-grouping-known-hosts` is also set)

---

## Summary

| Before | After |
|---|---|
| Inventory shape could be unstable across runs | Inventory shape is stable and version-controllable |
| No explicit role-based grouping | Groups declared explicitly by the operator |
| Single device family per run | Multiple device families merged into one inventory directory |
