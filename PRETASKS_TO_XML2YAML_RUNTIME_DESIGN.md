# Design Note: Move Playbook pre_tasks Logic into jtaf-xml2yaml Runtime Outputs

## Purpose

Define a design that removes heavy runtime pre_tasks from provider playbooks by shifting variable layering and merge resolution into jtaf-xml2yaml at generation time.

This document is design only. No implementation is included.

## Problem Statement

Current provider playbooks for VQFX and SRX contain pre_tasks that:
1. Load layered data from a central hierarchy (group_vars/all, group_vars/<group>/all, host_vars/<host>).
2. Build merged runtime facts.
3. Pass those merged facts to role templates.

This makes playbooks heavy and fragile because runtime behavior depends on custom merge orchestration.

## Target Outcome

1. Keep provider playbooks minimal (hosts, connection, role call).
2. Make jtaf-xml2yaml emit runtime-ready host variables per provider workspace.
3. Preserve canonical DRY hierarchy in examples/ansible/ansible_files for human review and future regeneration.
4. Keep XML rendering parity for VQFX and SRX.

## Non-Goals

1. No changes to template semantic intent in this design.
2. No removal of canonical ansible_files hierarchy.
3. No immediate rewrite of every filter plugin behavior unless required by runtime compiled payload format.

## Current vs Proposed Ownership

Current ownership:
- Runtime layering: playbook pre_tasks + role tasks + filters.
- Generation: jtaf-xml2yaml writes canonical hierarchy.

Proposed ownership:
- Runtime layering: jtaf-xml2yaml (compile effective host payloads).
- Playbook: native Ansible var loading only.
- Role tasks: minimal normalization and template render.

## High-Level Design

### A) Dual Output Model in jtaf-xml2yaml

jtaf-xml2yaml writes two views:

1. Canonical DRY view (existing behavior retained):
- group_vars/all.yml
- group_vars/<type>/all.yml
- host_vars/<host>.yaml deltas
- inventory groups

2. Runtime compiled view (new):
- Per-host fully resolved payload files for provider playbook execution.
- Output located under provider role workspace (or explicit runtime output path).

### B) Runtime Payload Contract

For each managed host in the active run:
1. Compute effective payload using the same merge semantics used for parity work.
2. Apply merge directives and list directives during generation.
3. Strip meta keys before writing runtime payload.
4. Persist as host_vars/<host>.yml in provider runtime output directory.

These runtime compiled host_vars are intended to be the values that Ansible loads for each inventory host and supplies to the corresponding provider role during template rendering. The role may still normalize them once into jtaf_effective, but it should not need to reconstruct the hierarchy at runtime.

Result: Playbook does not need layered include_vars/set_fact pre_tasks.

### C) Runtime host_vars Lifecycle and Storage Model

1. Runtime compiled host_vars are stored on disk as generated files, not kept only as transient in-memory data.
2. They are written into the provider-local playbook workspace, typically under host_vars/ for that role.
3. They are persistent execution artifacts: they remain present between playbook runs until the next regeneration, cleanup, or explicit reset.
4. They are not the canonical DRY source of truth; examples/ansible/ansible_files remains the canonical generated hierarchy.
5. They are derived artifacts and may be overwritten on subsequent jtaf-xml2yaml conversion runs.
6. They should be treated as generated runtime inputs for rendering, not as hand-maintained files.

### D) Playbook Simplification

Provider playbooks become minimal and stable:
- no jtaf_vars_root
- no layered include_vars
- no pre_tasks merge orchestration

## File-by-File Changes Required (Design)

## 1) junosterraform/jtaf-xml2yaml

Required design changes:
1. Add a single runtime-output-dir for compiled host_vars emission.
2. Add compile phase after canonical hierarchy computation:
- Resolve per-host effective data from global + per-type + host levels.
- Resolve merge directives/list directives.
- Remove meta keys from runtime outputs.
3. Write runtime compiled host_vars to the provider workspace runtime-output-dir.
4. Keep existing canonical output behavior unchanged.
5. Ensure repeated runs are merge-safe and do not clobber foreign provider data.

Potential CLI addition:
1. --runtime-output-dir <path>

runtime-output-dir is the provider role workspace path already supplied for that role, for example:
- examples/ansible/ansible-provider-junos-vqfx-ansible-role
- examples/ansible/ansible-provider-junos-srx-ansible-role

Acceptance criteria:
1. Canonical hierarchy is still generated as before.
2. Runtime compiled host_vars are complete enough for template render without pre_tasks.

## 2) examples/ansible/convert-ansible.sh

Required design changes:
1. Keep current generation of canonical ansible_files.
2. Add generation step(s) that produce runtime compiled vars for:
- ansible-provider-junos-vqfx-ansible-role
- ansible-provider-junos-srx-ansible-role
3. Ensure both provider workspaces receive runtime compiled host_vars for their devices under the provider runtime-output-dir.
4. Keep reset/cleanup behavior deterministic for both canonical and runtime outputs.

Acceptance criteria:
1. Single conversion flow prepares both canonical and runtime execution views.
2. No manual copying from ansible_files into provider directories.

## 3) examples/ansible/ansible-provider-junos-vqfx-ansible-role/jtaf-playbook.yml

Required design changes:
1. Remove heavy pre_tasks block.
2. Remove jtaf_vars_root and file lookup logic.
3. Keep minimal play structure only.

Acceptance criteria:
1. Playbook renders with runtime host_vars generated by jtaf-xml2yaml.

## 4) examples/ansible/ansible-provider-junos-srx-ansible-role/jtaf-playbook.yml

Required design changes:
1. Remove heavy pre_tasks block.
2. Remove jtaf_vars_root and file lookup logic.
3. Keep minimal play structure only.

Acceptance criteria:
1. Playbook renders with runtime host_vars generated by jtaf-xml2yaml.

## 5) examples/ansible/ansible-provider-junos-vqfx-ansible-role/roles/vqfx-ansible-role_role/tasks/main.yml

Required design changes:
1. Remove runtime layered combine assumptions tied to pre_tasks facts.
2. Build jtaf_effective directly from loaded runtime vars (already compiled).
3. Keep template render path stable.

Acceptance criteria:
1. No dependency on jtaf_group_all/jtaf_group_merged/jtaf_host_vars facts from pre_tasks.

## 6) examples/ansible/ansible-provider-junos-srx-ansible-role/roles/srx-ansible-role_role/tasks/main.yml

Required design changes:
1. Same simplification as VQFX role.
2. Keep render behavior stable for SRX template expectations.

Acceptance criteria:
1. No dependency on pre_tasks-created layered facts.

## 7) examples/ansible/ansible-provider-junos-vqfx-ansible-role/filter_plugins/jtaf_filters.py

Required design review:
1. Verify if runtime compiled payloads make merge filters optional at runtime.
2. Keep plugin for compatibility until migration is complete.
3. Optionally narrow plugin scope to normalization-only operations if no longer needed for merge orchestration.

## 8) examples/ansible/ansible-provider-junos-srx-ansible-role/filter_plugins/jtaf_filters.py

Required design review:
1. Same as VQFX plugin review.
2. Retain compatibility during transition.

## 9) junosterraform/jtaf-ansible

Required design changes:
1. Ensure generated playbook remains minimal by default.
2. Ensure generated task skeleton assumes runtime vars are already prepared by jtaf-xml2yaml.
3. Document generated role expectations clearly in comments.

Acceptance criteria:
1. Newly generated provider roles do not reintroduce heavy pre_tasks patterns.

## 10) README.md

Required documentation changes:
1. Document dual output model (canonical + runtime compiled).
2. Document new runtime output flags/usage.
3. Document ownership boundary: merge/layering happens in generator, not playbook.
4. Update examples for convert flow and execution flow.

## 11) PR_CHANGE_SUMMARY.md

Required documentation changes:
1. Add section describing migration from runtime pre_tasks to generator-compiled runtime vars.
2. Track file-level updates once implementation is done.

## 12) Optional: XML_SYNC_REAPPLY_TRACKER_2026-04-29.md

Required documentation changes:
1. Add note that parity now depends on generator compile phase, not playbook pre_tasks.
2. Keep historical context for regression tracking.

## Data and Merge Semantics (Design Contract)

Compile logic contract in jtaf-xml2yaml:
1. Effective(host) = DeepMerge(global_all, group_type_all, host_delta)
2. list_merge behavior must preserve prior parity behavior for list of dicts.
3. _merge_directive and _merge_list_directives are resolved before runtime file write.
4. Meta keys are excluded from runtime payload artifacts.

## Migration Plan (No-Code Design)

Phase 1: Add runtime compile output in jtaf-xml2yaml.
Phase 2: Update convert-ansible.sh to emit canonical + runtime outputs.
Phase 3: Simplify VQFX/SRX playbooks and role tasks.
Phase 4: Validate XML parity for all VQFX and SRX hosts.
Phase 5: Update docs and regression tracker.

## Validation Plan

1. Render VQFX generated XML files and compare against examples/evpn-vxlan-dc/dc1 originals.
2. Render SRX generated XML files and compare against examples/evpn-vxlan-dc/dc1 originals.
3. Confirm missing/extra leaf-value checks are zero (excluding known metadata keys such as JTAF_ANSIBLE).
4. Confirm playbooks run without pre_tasks merge logic.
5. Confirm repeated conversion runs remain stable.

## Risks and Mitigations

1. Risk: Runtime and canonical outputs diverge over time.
- Mitigation: Build runtime from canonical computed structures in one execution path.

2. Risk: List merge parity regressions.
- Mitigation: Reuse existing merge directive semantics in generator compile phase and run full XML parity comparison.

3. Risk: Backward compatibility for teams depending on current pre_tasks.
- Mitigation: Keep compatibility window where filter plugins remain available and document migration clearly.

## Rollback Strategy

If migration causes regressions:
1. Temporarily keep simplified runtime output disabled via a feature flag.
2. Re-enable existing pre_tasks path for affected provider playbook.
3. Use parity comparison reports to isolate merge semantic mismatch.

## Implementation Readiness Checklist

1. Design approved.
2. CLI contract for runtime output finalized.
3. Conversion script flow agreed.
4. Playbook simplification agreed for both VQFX and SRX.
5. Validation corpus and comparison script ready.
6. Documentation update plan approved.
