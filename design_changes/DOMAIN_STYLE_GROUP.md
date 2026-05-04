# Domain-Style Grouping Changes

This document captures the recent changes made to domain-style grouping in `jtaf-xml2yaml` and how those changes affect generated Ansible output.

## Goal

Generate a practical, deterministic inventory/group_vars layout similar to the curated target structure in `ansible_files_domain_target`, while keeping generic fallback behavior for non-domain datasets.

## What Was Added

### 1) Deterministic domain-style grouping

A new grouping pass was added to prefer role-oriented groups when hostnames/payloads allow reliable classification.

Main function:
- `derive_domain_style_groups(host_var_entries, min_hosts=2)`

Supporting helpers:
- `detect_role_label(hostname)`
- `is_firewall_like_host(hostname, payload)`

Domain-style grouping shape:
- `group1`: non-firewall aggregate parent (when needed)
- `group2..N`: role child groups in fixed order: borderleaf, leaf, spine
- final group: firewall group (for firewall-like hosts)

Determinism is enforced by:
- fixed role ordering
- sorted unique host lists
- stable group numbering sequence

### 2) Integration into the existing common-group pipeline

`derive_common_host_groups(...)` now calls `derive_domain_style_groups(...)` first.

Behavior:
- if domain groups are detected, they are returned first (with stable numeric ordering and max-count cap)
- if not detected, logic falls back to the previous generic overlap/scoring-based grouping algorithm

This keeps broad applicability while strongly preferring the desired domain layout where possible.

### 3) Inventory parent-group host membership

Inventory generation was updated so ancestor groups also receive host membership (not only leaf groups).

Main function impacted:
- `build_hierarchical_inventory_groups(parent_map, host_to_leaf_group)`

Result:
- parent groups such as `group1` include explicit hosts
- child relationships are still emitted via `[group1:children]`

### 4) group_vars default path alignment

Default output directory for group vars was standardized to:
- `<output>/group_vars`

(Previously some paths/messages referenced `group-vars`.)

## Output Characteristics After Change

For the dc1 EVPN/VXLAN sample domain, generated layout now follows the target pattern:
- `group1` = borderleaf + leaf + spine hosts
- `group2` = borderleaf hosts
- `group3` = leaf hosts
- `group4` = spine hosts
- `group5` = firewall hosts

Hosts file includes both:
- explicit host sections for each group (including `group1`)
- parent-child linkage sections such as `[group1:children]`

## Validation Summary

The changes were validated with:
- unit/regression tests for hierarchical grouping logic
- script coverage tests
- regenerated `examples/ansible/ansible_files` output
- end-to-end Ansible playbook runs:
  - vqfx role with `--limit group1`
  - srx role with `--limit group5`

All validation passed (only existing Ansible deprecation warnings remained).

## Notes

- Existing tunables (`max depth`, `max count`, `min benefit`, `min new paths`) are still honored.
- Single-host leaf group collapsing and conflict fallback logic remain active.
- Domain-style grouping is preferred, not forced; generic fallback remains available for other datasets.
