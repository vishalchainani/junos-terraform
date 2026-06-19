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

## Simple Flat Example — Switches

Suppose you have nine QFX devices across two data-centre pods. Create a file, e.g. `qfx.grouping.hosts`:

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

Then run:

```bash
jtaf-xml2yaml \
  -j ansible-provider-junos-vqfx-evpn-vxlan/trimmed_schema.json \
  -x examples/evpn-vxlan-dc/dc1/*{spine,leaf}*.xml \
     examples/evpn-vxlan-dc/dc2/*spine*.xml \
  -d ansible-evpn-vxlan-deploy \
  --grouping-hosts-file ansible-evpn-vxlan-deploy/qfx.grouping.hosts
```

The tool writes:

```
ansible-evpn-vxlan-deploy/
├── inventory.ini
├── group_vars/
│   ├── all.yaml              ← config keys shared by ALL nine devices
│   ├── borderleaf/
│   │   └── all.yaml          ← keys shared only by dc1-borderleaf1 and dc1-borderleaf2
│   ├── leaf/
│   │   └── all.yaml          ← keys shared only by the three leaf devices
│   └── spine/
│       └── all.yaml          ← keys shared only by the four spine devices
└── host_vars/
    ├── dc1-borderleaf1.yaml  ← unique delta for this device
    ├── dc1-leaf1.yaml
    └── ...
```

The generated `inventory.ini` will contain the `[borderleaf]`, `[leaf]`, and `[spine]` group sections exactly as declared in the grouping file.

---

## Multiple Device Families in One Inventory Directory

A key design point is that two separate `jtaf-xml2yaml` runs — each with its own `--grouping-hosts-file` — can target the **same `-d` output directory**. The tool merges inventory sections without clobbering existing groups.

For example, to add SRX firewall devices alongside the QFX switches:

Create `firewall.grouping.hosts`:

```ini
[all]
dc1-firewall1
dc1-firewall2
dc2-firewall1
dc2-firewall2

[firewall]
dc1-firewall1
dc1-firewall2
dc2-firewall1
dc2-firewall2
```

Run the second conversion targeting the same directory:

```bash
jtaf-xml2yaml \
  -j ansible-provider-junos-srx-ansible-role/trimmed_schema.json \
  -x examples/evpn-vxlan-dc/dc1/dc1-*firewall*.xml \
     examples/evpn-vxlan-dc/dc2/dc2-*firewall*.xml \
  -d ansible-evpn-vxlan-deploy \
  --grouping-hosts-file ansible-evpn-vxlan-deploy/firewall.grouping.hosts
```

After both runs, `ansible-evpn-vxlan-deploy/inventory.ini` contains all groups from both files:

```ini
[all]
dc1-borderleaf1
dc1-borderleaf2
dc1-leaf1
...
dc2-firewall2

[borderleaf]
dc1-borderleaf1
dc1-borderleaf2

[leaf]
...

[firewall]
dc1-firewall1
dc1-firewall2
dc2-firewall1
dc2-firewall2
```

And `group_vars/firewall/all.yaml` only contains config shared among the four firewall devices — completely separate from the QFX group vars — even though everything lives under one output directory.

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
