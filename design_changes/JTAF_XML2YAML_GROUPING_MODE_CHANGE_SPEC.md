# jtaf-xml2yaml Grouping Mode Change Spec (No Code Changes Yet)

## Objective
Implement a user-defined grouping workflow in `jtaf-xml2yaml`:

1. Remove `auto` grouping mode completely.
2. Default behavior is user grouping only.
3. Script runs non-interactively by default and requires `grouping.hosts` via CLI argument.
4. Script validates `grouping.hosts` and reports syntax errors.
4. Remove automatic group discovery and numeric group naming/numbering scheme.

This document describes **what to change** and anticipated caveats. It does not implement code.

---

## Current State (Relevant Parts)
`junosterraform/jtaf-xml2yaml` currently includes:

- Automatic grouping logic:
  - `derive_common_host_groups(...)`
  - `build_group_parent_map(...)`
  - `collapse_daisy_chains(...)`
  - `resolve_group_selection_settings(...)`
  - `build_shape_report(...)`
- Numeric group allocation/renaming logic:
  - `allocate_group_range(...)`
  - `renumber_groups(...)`
  - `inventory_group_sort_key(...)` with numeric preference
- Hierarchy payload writing:
  - `write_group_hierarchy_payloads(...)`
  - `rewrite_host_vars_with_hierarchy(...)`
- Registry/global shared payload path:
  - `load_shared_all_state(...)`, `write_shared_all_state(...)`, provider registry metadata in `group_vars/all.yaml`

---

## Required Design Changes

## 1) Remove Auto Mode, Keep User Mode Only

### Behavior
At runtime, grouping mode should no longer be selected. The script must always run in user-defined grouping mode.

### Parser/API updates
Update CLI and parser behavior:
- Remove `--grouping-mode` option entirely.
- Make non-interactive behavior the default for all invocations.
- Require `--grouping-hosts-file <path>` for both human and CI usage.
- Optional: add `--interactive` as an explicit opt-in fallback mode.

---

## 2) Mandatory Argument Input + Validation

### Default flow (non-interactive)
1. User passes `--grouping-hosts-file <path>`.
2. Script validates file.
3. If syntax errors exist, print all errors and exit non-zero.

### Optional interactive fallback (only if explicitly enabled)
If `--interactive` is provided:
1. Print sample `grouping.hosts` file before path prompt.
2. Prompt:
  - `Enter path to grouping.hosts:`
3. Validate file.
4. If syntax errors exist, print all errors and re-prompt.

### Sample to display in prompt
Use this exact sample (Ansible INI-style inventory):

```ini
[all]
dc1-borderleaf1
dc1-borderleaf2
dc1-firewall1
dc1-firewall2
dc1-leaf1
dc1-leaf2
dc1-leaf3
dc1-spine1
dc1-spine2

[switch]
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

[firewall]
dc1-firewall1
dc1-firewall2

[switch:children]
borderleaf
leaf
spine
```

### Validation scope (must-have)
Validate and report line-numbered errors:

- File exists and readable.
- Valid section headers (`[name]` or `[name:children]` only).
- No malformed section headers.
- Host lines only appear under host sections.
- Child-group lines only appear under `:children` sections.
- No duplicate section name declarations.
- No undefined child group references.
- No cyclic group ancestry (`A -> B -> A`).
- Optional strict check: all hosts in groups must be known hosts from XML input set.

### Existing function reuse
- Extend or wrap `load_inventory_sections(...)` so it can return structured parse errors.
- Add new `validate_grouping_hosts(...)` helper returning `(ok, errors, parsed_inventory)`.

---

## 3) Remove Auto Execution Path Entirely

### Required outcome
There is no auto mode path at all. Any code branch that previously executed auto behavior must be removed.

### Implementation notes
- Keep only user-defined hierarchy execution paths.
- Always derive `group_vars/<group>/all.yaml` and `host_vars/<host>.yaml` from validated `grouping.hosts` input.
- Fail fast if grouping file is missing/invalid.

---

## 4) Remove Automatic Group Discovery + Numbering

### Remove (or hard-disable) these paths
- `derive_common_host_groups(...)`
- `build_group_parent_map(...)`
- `collapse_daisy_chains(...)`
- `resolve_group_selection_settings(...)`
- `build_shape_report(...)`
- `allocate_group_range(...)`
- `renumber_groups(...)`
- Any code that emits `group1`, `group2`, etc.
- Numeric-first sorting behavior in `inventory_group_sort_key(...)`.

### Keep/repurpose
- Keep inventory parser/writer helpers, but repurpose for user-defined groups only.
- Keep payload merge/subtract utilities; they are still useful in `user` mode when computing group deltas and host deltas.

---

## User Mode Data Flow (Target)

1. Parse XML -> full per-host payloads.
2. Read validated `grouping.hosts` from `--grouping-hosts-file`.
3. Build parent map and group membership from user file (no auto-discovery).
4. Compute shared payload per declared group from declared hosts.
5. Write:
   - `group_vars/all.yaml` (only if desired in user mode)
   - `group_vars/<group>/all.yaml` for declared groups
6. Rewrite `host_vars/<host>.yaml` as strict deltas against ancestry baseline.
7. Write resulting inventory preserving user group names and structure.

Note: In user mode, group names must remain exactly as provided. No renumbering.

---

## `convert-ansible.sh` Integration (Recommended)

Update `examples/ansible/convert-ansible.sh` to always pass explicit values for deterministic runs:

- `--grouping-hosts-file <path>` (required)

This ensures identical behavior for humans and CI/CD.

---

## Caveats / Risks

- Human users must now always provide `--grouping-hosts-file`; this is a deliberate UX tradeoff for consistency.
- Existing registry metadata in `group_vars/all.yaml` may require migration cleanup when removing mode branching.
- Multi-provider runs currently depend on shared global payload logic; behavior must remain consistent under user-only mode.
- Validation strictness for unknown hosts may block partial runs (e.g., when only a subset of devices is converted in one call).
- If user-defined groups overlap heavily, host-delta resolution may still require conflict handling similar to current sibling conflict fallback.

---

## Migration/Compatibility Strategy

1. Add release note: auto mode removed; grouping is user-defined only.
2. Runtime behavior:
  - Default (human + CI): require `--grouping-hosts-file`; fail fast if missing.
  - Optional fallback: `--interactive` can prompt for path if enabled.
3. On first run after upgrade, clear legacy `group1/...` structures.

---

## Testing Checklist

- `user` mode happy path with valid grouping file.
- `user` mode malformed section syntax.
- `user` mode undefined child group.
- `user` mode cycle detection.
- `user` mode unknown host handling.
- Missing `--grouping-hosts-file` fails fast by default.
- No `group1`, `group2`, etc. generated in any mode.
- Re-run stability (idempotence) across repeated executions.
- `convert-ansible.sh` works for both vqfx and srx calls with same grouping source.

---

## Proposed New CLI Options (for implementation phase)

- `--grouping-hosts-file <path>`
- `--interactive` (optional fallback prompt mode)

Behavior:
- Default: if `--grouping-hosts-file` is missing, exit with actionable error.
- If `--interactive` is provided, prompt for path and validate.

---

## Summary
The code should move from **auto-discovered numeric grouping** to **user-defined grouping only**:

- `grouping.hosts` is mandatory input via CLI argument by default.
- Groups and hierarchy are fully defined by validated user input.

This matches your requirement for user-driven naming and structure while eliminating automatic grouping and numbering.
