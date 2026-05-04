# XML Sync Re-Apply Tracker (2026-04-29)

## Purpose
This document tracks the re-application of XML parity fixes for VQFX and SRX Ansible render pipelines after regression.

## Requested Outcome
- Re-apply hardened merge and template changes.
- Re-render configurations.
- Re-compare generated XML against original device XML.
- Keep a file-level audit trail of what changed.

## Files Updated

### VQFX provider changes
1. examples/ansible/ansible-provider-junos-vqfx-ansible-role/jtaf-playbook.yml
- Added explicit layered variable loading pre_tasks.
- Added jtaf_vars_root and facts:
  - jtaf_group_all
  - jtaf_group_merged
  - jtaf_host_vars
- Added file-based loading from examples/ansible/ansible_files:
  - group_vars/all.yaml
  - group_vars/<group>/all.yaml
  - host_vars/<host>.y*ml

2. examples/ansible/ansible-provider-junos-vqfx-ansible-role/roles/vqfx-ansible-role_role/tasks/main.yml
- Replaced vars dictionary scan approach with deterministic layered combine:
  - jtaf_group_all
  - combine with jtaf_group_merged (recursive=True, list_merge='append')
  - combine with jtaf_host_vars (recursive=True, list_merge='append')
- Updated directive application step to include jtaf_remove_meta.

3. examples/ansible/ansible-provider-junos-vqfx-ansible-role/filter_plugins/jtaf_filters.py
- Added deep merge helper for nested dict/list merge behavior.
- Added list merge-by-key capability:
  - _merge_list_by_key
  - _apply_list_directives
- Updated directive processor to apply _merge_list_directives recursively.
- Updated metadata cleanup to remove _applied_directive plus merge metadata keys.

4. examples/ansible/ansible-provider-junos-vqfx-ansible-role/roles/vqfx-ansible-role_role/templates/template.j2
- Fixed policy term protocol rendering.
- Added iterable handling so protocol lists render as repeated protocol XML entries.

### SRX provider changes
5. examples/ansible/ansible-provider-junos-srx-ansible-role/jtaf-playbook.yml
- Added explicit layered variable loading pre_tasks (same pattern as VQFX).
- Added jtaf_vars_root and facts for layered combine.

6. examples/ansible/ansible-provider-junos-srx-ansible-role/roles/srx-ansible-role_role/tasks/main.yml
- Replaced vars dictionary scan with explicit layered combine.
- Added jtaf_remove_meta to post-directive cleanup.

7. examples/ansible/ansible-provider-junos-srx-ansible-role/filter_plugins/jtaf_filters.py
- Added deep recursive merge behavior.
- Added list merge-by-key support and directive application.
- Updated metadata removal behavior as in VQFX filter plugin.

8. examples/ansible/ansible-provider-junos-srx-ansible-role/roles/srx-ansible-role_role/templates/template.j2
- Added version element rendering.
- Added syn-flood queue-size rendering from undocumented.queue_size.
- Added gRPC undocumented clear-text address and port rendering.
- Added skip-authentication rendering under gRPC undocumented block.

### Shared vars change
9. examples/ansible/ansible_files/group_vars/all.yaml
- Restored gRPC clear_text defaults under:
  - system.services.extension_service.request_response.grpc.undocumented.clear_text.address
  - system.services.extension_service.request_response.grpc.undocumented.clear_text.port

## Validation Performed

### Syntax and diagnostics
- Checked edited files for diagnostics.
- Result: no errors reported.

### Render runs
1. VQFX full render run
- Playbook: examples/ansible/ansible-provider-junos-vqfx-ansible-role/jtaf-playbook.yml
- Inventory: examples/ansible/ansible_files/hosts
- Result: successful; all dc1 VQFX-related files regenerated.

2. SRX firewall render run
- Playbook: examples/ansible/ansible-provider-junos-srx-ansible-role/jtaf-playbook.yml
- Inventory: examples/ansible/ansible_files/hosts
- Limit: dc1-firewall1,dc1-firewall2
- Result: successful; both firewall configs regenerated.

### XML parity comparison results
Compared generated files against originals in examples/evpn-vxlan-dc/dc1.

#### VQFX
- dc1-borderleaf1.xml: MISSING 0, EXTRA 0
- dc1-borderleaf2.xml: MISSING 0, EXTRA 0
- dc1-leaf1.xml: MISSING 0, EXTRA 0
- dc1-leaf2.xml: MISSING 0, EXTRA 0
- dc1-leaf3.xml: MISSING 0, EXTRA 0
- dc1-spine1.xml: MISSING 0, EXTRA 0
- dc1-spine2.xml: MISSING 0, EXTRA 0

#### SRX
- dc1-firewall1.xml: MISSING 0, EXTRA 0
- dc1-firewall2.xml: MISSING 0, EXTRA 0

## Final Status
- Regression fixes have been re-applied.
- Current generated VQFX and SRX firewall XML files are in sync with original dc1 XML files (no missing or extra tag/value elements in the comparison method used).