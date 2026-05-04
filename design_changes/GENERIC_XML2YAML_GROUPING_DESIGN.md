# Generic Hierarchical Grouping Design for jtaf-xml2yaml

## Purpose

Define a generic, provider-scoped grouping design for jtaf-xml2yaml that directionally reproduces the output style in examples/ansible/ansible_files_target without hard-coded domain knowledge.

This document is design only. No implementation is included.

## Design Goals

1. Recreate as closely as practical the hosts and YAML structure style in examples/ansible/ansible_files_target.
2. Remove all built-in domain assumptions and hostname semantics from grouping decisions.
3. Add tunable controls for group depth, breadth, and merge-directive behavior.
4. Keep output directionally similar, not byte-identical.
5. Keep algorithm generic for any network and provider schema.
6. Use both XML payloads and provider schema metadata where useful.
7. Preserve files and inventory entries managed by other providers, but exclude them from grouping calculations for the active provider.
8. Continue supporting repeated runs across different providers with safe hosts file updates.
9. Never rely on hostname parsing for grouping.

## Current Gaps to Address

The current flow still includes domain-specific pathing in some logic and test expectations. The design replaces that with schema and payload driven grouping only.

## Invariants

1. Group derivation uses payload structure and values only.
2. Hostname content is never used to infer group semantics.
3. Provider scoping is explicit and deterministic.
4. Group names are generic and stable for a given input set.
5. Merge directives are optional, deterministic, and tunable.
6. Existing non-active-provider data is preserved.
7. group_vars/all.yaml is a shared file across providers and is never replaced wholesale by one provider run.

## Proposed High-Level Pipeline

### Stage 1: Provider Scope Resolution

Inputs:
- Trimmed schema path
- XML files supplied to this run

Provider key resolution order:
1. Derive from the Ansible role name contained inside the provider.
2. Normalize role name to a stable key token.
3. If role name cannot be resolved, fail fast with a clear diagnostic.

Provider-key source rule:
- provider-key MUST come from the role directory name under provider roles/.
- Example: roles/srx-ansible-role_role -> provider-key srx-ansible-role_role.
- User override is not supported.

Managed host set for this run:
- Hosts parsed from XML files in this invocation.

Tracked state:
- Store registry metadata in a meta section inside group_vars/all.yaml.
- Registry stores provider key, managed hosts, and group allocation history.
- Meta registry keys are reserved and excluded from grouping calculations.
- group_vars/all.yaml is shared across all providers writing into the same output directory.

Registry placement in group_vars/all.yaml (design intent):

```yaml
meta:
	jtaf_registry:
		version: 1
		providers:
			<provider_key>:
				managed_hosts: []
				group_range:
					start: 1
					end: 5
```

Shared all.yaml contract:
1. group_vars/all.yaml is a shared artifact for all providers in the output directory.
2. The active provider updates only:
- its own registry branch under meta.jtaf_registry.providers.<provider_key>
- shared non-meta payload using a merge policy that preserves foreign-provider state
3. Full-file overwrite is forbidden.
4. meta keys are excluded from grouping/scoring/delta math.

### Stage 2: Data Loading and Partitioning

1. Load full payload cache.
2. Load existing group_vars/all.yaml and split into meta and non-meta sections.
3. Partition hosts into:
- active provider hosts
- non-active-provider hosts
4. Run calculations only on active provider hosts.
5. Keep non-active-provider host vars and inventory sections unchanged unless explicitly requested.

### Stage 3: Generic Candidate Group Discovery

Candidate discovery uses only payload-derived signals:
1. Shared leaf path and value support sets.
2. Intersections of strong support sets.
3. Optional structural signature similarity using path fingerprints.
4. Optional schema-aware normalization for typed list handling.

No role names, no hostname prefixes, no product family labels are used.

### Stage 4: Candidate Scoring and Selection

Score each candidate using tunable components:
- Reuse benefit score from shared leaf count times host count gain.
- Novelty score from new paths not already covered by selected supersets.
- Compactness penalty for near-duplicate subsets.
- Optional schema-weighted score for high-signal paths.

Select candidates with constraints:
- minimum hosts per group
- maximum groups
- minimum benefit
- minimum new paths

### Stage 5: Hierarchy Construction

Build a nearest-superset parent map with controls:
- max depth
- max children per node
- minimum parent overlap ratio
- maximum single-child chain length

Anti-daisy-chain guardrail:
1. Detect any path where each group has exactly one child for more than the configured chain limit.
2. Collapse the chain by skipping intermediate parents with low unique contribution.
3. Preserve data semantics by recomputing strict deltas after collapse.
4. If collapse cannot maintain deterministic deltas, fail fast with diagnostics.

Chain collapse heuristic:
- An intermediate node is removable when its unique leaf-path contribution is below a configured threshold.
- Parent/child links are rewired to the nearest remaining ancestor.
- Host-to-leaf assignment is recomputed after rewiring.

Then assign dense generic group names by deterministic order:
- breadth-first from roots
- then by descending host count
- then lexical tie-break

Resulting names:
- group1, group2, group3, and so on

Dense numbering is local to the active provider allocation strategy.

### Stage 6: Payload Hoisting and Delta Writing

1. Compute active-provider shared payload candidate and merge it into shared group_vars/all.yaml using non-destructive update rules.
2. Compute full group payloads from member intersections.
3. Convert group payloads into strict deltas against ancestry.
4. Rewrite host payloads as strict deltas against global plus ancestry baseline.

Conflict handling:
- detect sibling conflicts on overlapping inherited paths
- remove ambiguous sibling paths
- emit host-specific fallback deltas

### Stage 7: Merge Directive Annotation

Directive strategy is configurable:
- off
- minimal
- balanced
- aggressive

Base rules:
- dict deltas can annotate merge_recursive
- list behavior can be replace, union, append, or merge_by_key
- merge_by_key key inference can use schema hints then fallback key preferences

### Stage 8: Inventory Update Strategy

Inventory update is provider-aware:
1. Keep all existing host entries.
2. Update only sections tied to active provider managed groups.
3. Preserve unknown sections and unrelated provider sections.

To avoid collisions across providers with generic group names, always use Mode A:

Mode A, mandatory for multi-provider safety:
- provider-scoped numeric namespace allocation via registry.
- example active provider may get group6 to group10 if group1 to group5 already allocated elsewhere.

## New Tuning Parameters

Existing parameters remain supported. Add the following:

1. group-selection-profile
- values: compact, balanced, expansive
- pre-sets for thresholds and breadth/depth defaults.

2. common-host-groups-min-hosts
- minimum group size.

3. common-host-groups-max-children
- max children for any parent group.

4. common-host-groups-max-single-child-chain
- maximum allowed consecutive single-child depth.

5. common-host-groups-chain-collapse-min-unique-paths
- minimum unique leaf paths required to keep an intermediate node in a single-child chain.

6. common-host-groups-min-parent-overlap
- minimum strict overlap ratio required for parent-child linkage.

7. merge-directive-profile
- values: off, minimal, balanced, aggressive.

8. merge-list-strategy-default
- values: replace, union, append, merge_by_key_auto.

9. merge-list-key-preference
- ordered list such as name,id,key,uuid.

10. preserve-foreign-provider-data
- default true.
- preserve host vars and hosts sections not owned by active provider.

11. dry-run-shape-report
- emit shape metrics without writing files.

## Provider Scoping Rules

1. Calculations include only active provider managed hosts.
2. Non-active-provider files remain untouched.
3. Inventory sections not owned by active provider remain untouched.
4. Ownership mapping is tracked in group_vars/all.yaml under meta.jtaf_registry.
5. meta and meta.jtaf_registry are excluded from shared-path collection, scoring, and delta subtraction logic.
6. group_vars/all.yaml is shared across providers; active provider writes must be additive/merge-safe and must not drop foreign-provider metadata.

## Directional Similarity Targets

Use structural metrics to evaluate closeness to target-style output:

1. Group depth distribution similar to target profile.
2. Parent-child fanout similar to target profile.
3. Host vars size reduced versus flat per-host payloads.
4. Group vars carry majority of reusable payloads.
5. Merge directives present where repeated nested list and dict merges occur.

These are directional metrics, not exact file matching constraints.

## Suggested Internal Refactor Boundaries

1. Group derivation module
- generic candidate generation and scoring only.

2. Provider scope module
- host partitioning, ownership registry, preservation logic.

3. Hierarchy shaping module
- parent mapping, numbering, depth and breadth constraints.

4. Delta and directives module
- group delta generation, host delta rewrite, directive annotation.

5. Inventory merge module
- provider-aware section updates and collision-safe numbering.

## Test Plan

### Unit Tests

1. No hostname semantic dependency
- grouping unchanged when hostnames are randomized.

2. Generic grouping determinism
- same input produces same group memberships and numbering.

3. Provider scope isolation
- active provider calculations ignore foreign provider hosts.

4. Inventory preservation
- unrelated groups stay unchanged after active provider run.

5. Merge directive profiles
- off, minimal, balanced, aggressive produce expected annotation levels.

6. Anti-daisy-chain guardrail
- no hierarchy path exceeds configured single-child chain length after collapse pass.

### Integration Tests

1. Run against examples/ansible data and compare shape metrics with ansible_files_target.
2. Run sequentially for two providers in same output directory and verify no destructive overwrite.
3. Re-run stability test to confirm deterministic outputs.

## Rollout

1. Implement the generic model as the only supported behavior.
2. Remove legacy code paths and compatibility flags.
3. Reject deprecated options with clear errors.

## Risks and Mitigations

1. Risk: under-grouping on small datasets.
- Mitigation: profile-based thresholds and minimum-host tuning.

2. Risk: over-grouping creating deep trees.
- Mitigation: max depth, max children, minimum novelty controls, and anti-chain collapse guardrail.

3. Risk: multi-provider group number collisions.
- Mitigation: provider ownership registry and shared numbering mode.

4. Risk: merge directive noise.
- Mitigation: directive profile and minimum-complexity threshold.

## Acceptance Criteria

1. No hostname-based grouping heuristics remain.
2. Active provider calculations exclude foreign provider files.
3. Foreign provider outputs are preserved.
4. Output hierarchy is directionally similar to ansible_files_target.
5. Group names are generic and dense within the chosen numbering mode.
6. Multi-provider sequential runs are safe and deterministic.
7. Tests cover invariants, provider scoping, and hierarchy shape.
8. group_vars/all.yaml remains a shared multi-provider artifact after sequential provider runs.
9. Hierarchy has no daisy-chain path beyond configured single-child chain limit.

## Implementation Todo List

Status legend:
- `not-started`
- `in-progress`
- `blocked`
- `done`

Update rules:
1. Set phase status before starting work.
2. Check boxes as items are completed.
3. Fill `Updated` and `Notes` each time status changes.

### Phase Tracker

| Phase | Status | Updated | Notes |
|---|---|---|---|
| Phase 0: Scope Lock and Baseline | done | 2026-04-20 | Baseline converter flow, target outputs, provider layout, and validation commands were mapped and exercised. |
| Phase 1: Provider Identity and Shared Registry | done | 2026-04-20 | Provider-key derivation, shared all.yaml registry persistence, and meta isolation are implemented. |
| Phase 2: Provider Scoping and Data Partitioning | done | 2026-04-20 | Active-provider grouping is isolated and foreign provider inventory and vars are preserved across runs. |
| Phase 3: Generic Candidate Discovery and Selection | done | 2026-04-20 | Generic payload-driven candidate selection replaced hostname and domain-specific grouping. |
| Phase 4: Hierarchy Shaping and Anti-Daisy-Chain | done | 2026-04-20 | Parent mapping, single-child chain collapse, and host leaf assignment are implemented and tested. |
| Phase 5: Delta Generation and Merge Directives | done | 2026-04-20 | Strict ancestry-aware deltas, conflict fallback, and merge directive hints are implemented. |
| Phase 6: Inventory and Numbering | done | 2026-04-20 | Mode A numbering, provider-owned range persistence, and sequential multi-provider inventory updates are working. |
| Phase 7: Documentation and Final Gates | in-progress | 2026-04-20 | Final lint and test gates are green; README and example consistency updates remain open. |

### Phase 0: Scope Lock and Baseline

Status: `done`

- [x] Freeze implementation scope against this document and acceptance criteria.
- [x] Map each design section to target code locations in jtaf-xml2yaml.
- [x] Prepare two-provider fixture set for shared-output validation.
- [x] Define canonical lint and test commands for repeatable execution.

Hygiene Cycle A:
- [x] Run baseline lint.
- [x] Run baseline targeted tests.
- [x] Run full test suite smoke.
- [x] Fix only environment or harness blockers.
- [x] Re-run until baseline is clean.

### Phase 1: Provider Identity and Shared Registry

Status: `done`

- [x] Implement provider-key derivation from provider role name only.
- [x] Remove provider-key override behavior and reject deprecated overrides.
- [x] Implement shared group_vars/all.yaml registry read and write under meta.jtaf_registry.
- [x] Implement strict separation of meta and non-meta content paths.
- [x] Ensure registry keys are excluded from grouping and scoring math.

Hygiene Cycle B:
- [x] Run lint on modified modules.
- [x] Run unit tests for provider-key derivation and registry persistence.
- [x] Fix failures.
- [x] Re-run lint and tests until green.

### Phase 2: Provider Scoping and Data Partitioning

Status: `done`

- [x] Implement active-provider versus non-active-provider host partitioning.
- [x] Restrict grouping calculations to active-provider hosts only.
- [x] Preserve non-active-provider host_vars and unrelated inventory sections.
- [x] Implement non-destructive shared all.yaml updates for active provider writes.

Hygiene Cycle C:
- [x] Run lint on partitioning and IO paths.
- [x] Run scoped tests for preservation and partition correctness.
- [x] Fix and re-run until stable.

### Phase 3: Generic Candidate Discovery and Selection

Status: `done`

- [x] Implement payload and schema-based candidate discovery.
- [x] Remove residual hostname and domain heuristics.
- [x] Implement benefit and novelty scoring with deterministic ordering.
- [x] Enforce selection constraints for size, count, and novelty.

Hygiene Cycle D:
- [x] Run lint on grouping and scoring modules.
- [x] Run unit tests for deterministic candidate selection.
- [x] Fix and re-run until all pass.

### Phase 4: Hierarchy Shaping and Anti-Daisy-Chain

Status: `done`

- [x] Implement nearest-superset parent mapping with depth, child-count, and overlap controls.
- [x] Implement single-child chain detection.
- [x] Implement chain collapse based on unique-path contribution threshold.
- [x] Rewire ancestry deterministically and recompute host-to-leaf assignments.
- [x] Fail fast when deterministic delta semantics cannot be preserved.

Hygiene Cycle E:
- [x] Run lint on hierarchy construction paths.
- [x] Run targeted tests for chain collapse and bounded chain depth.
- [x] Fix and re-run until constraints are enforced.

### Phase 5: Delta Generation and Merge Directives

Status: `done`

- [x] Implement strict group deltas against ancestry baseline.
- [x] Implement strict host deltas against global plus inherited baseline.
- [x] Implement sibling conflict fallback behavior.
- [x] Implement merge directive profiles and list strategy options.

Hygiene Cycle F:
- [x] Run lint on delta and directive modules.
- [x] Run unit tests for merge behavior and conflict handling.
- [x] Fix and re-run until deterministic outputs are stable.

### Phase 6: Inventory and Numbering

Status: `done`

- [x] Implement Mode A only numbering allocation from shared registry.
- [x] Enforce collision-safe provider-scoped group range ownership.
- [x] Implement provider-aware inventory updates that preserve unrelated sections.
- [x] Validate sequential multi-provider run stability.

Hygiene Cycle G:
- [x] Run lint on inventory and numbering logic.
- [x] Run integration tests for multi-provider sequential runs.
- [x] Fix and re-run until no destructive overwrite occurs.

### Phase 7: Documentation and Final Gates

Status: `in-progress`

- [ ] Update README and related docs to reflect generic-only behavior.
- [ ] Remove documentation of legacy or compatibility pathways.
- [ ] Verify examples are consistent with implemented behavior.

Hygiene Cycle H:
- [x] Run full lint across repository.
- [x] Run full test suite.
- [x] Fix residual issues.
- [x] Re-run full lint and tests until fully green.

### Final Release Checklist

- [ ] Confirm all acceptance criteria in this document are met.
- [x] Confirm shared group_vars/all.yaml behavior preserves foreign-provider data.
- [x] Confirm no hierarchy path exceeds configured single-child chain limit.
- [x] Confirm deterministic re-run stability for fixtures.
- [ ] Publish short implementation report with test evidence and residual risks.
