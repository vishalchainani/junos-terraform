# PR Change Summary

## Scope
- Updated Ansible variable generation to support hierarchical layering for mixed device types.
- Added optional delta mode to reduce repeated keys in per-device group vars.
- Updated example conversion script, tests, and docs accordingly.

## Core Behavior Changes
- jtaf-xml2yaml now supports Option B hierarchy:
  - group_vars/all.yml (global shared intersection)
  - group_vars/<type>/all.yml (per-type shared values)
  - host_vars/<host>.yaml (host-specific deltas)
- Inventory groups continue to be generated (e.g., vqfx, srx).
- Host deltas are computed from global + active device-type baseline.

## New Optional Flag
- Added --device-group-delta to jtaf-xml2yaml.
- In delta mode, group_vars/<type>/all.yml is written as delta against group_vars/all.yml.
- Empty per-type delta files are omitted, and empty <type> directories are removed.
- Default behavior remains unchanged when the flag is not provided.

## Script Update
- examples/ansible/convert-ansible.sh updated to pass:
  - --auto-detect-hierarchy
  - --device-group-delta

## Docs Updated
- README.md updated for Option B layout, repeated-run semantics, and new flag usage.
- HIERARCHICAL_GROUPS_WITH_DIRECTIVES.md updated for merge order, examples, CLI options, and delta-mode notes.
- jtaf-ansible generated task comments updated to reflect all -> <type> -> host precedence.

## Tests / Validation
- Updated tests in junosterraform/tests/test_hierarchical_groups.py.
- Updated tests in junosterraform/tests/test_workflow.py.
- Test results:
  - pytest junosterraform/tests/test_hierarchical_groups.py -q : 15 passed
  - pytest junosterraform/tests/test_workflow.py -q : 1 passed, 2 skipped

## Practical Outcome
- Better structure for mixed VQFX/SRX runs.
- Cleaner per-type files with optional delta mode.
- Backward-safe default behavior retained.