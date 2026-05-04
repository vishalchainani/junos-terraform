# Hierarchical Group-Var Design (Hierarchy-Only)

## Scope

This design defines a hierarchy-only variable model for generated YAML.

- No legacy flat group-var format.
- No backward compatibility mode.
- Deterministic merge order.
- Generic group naming and behavior (no implementation-type assumptions).

## Output Layout

Generated structure:

- `group-vars/all.yaml` (global baseline)
- `group-vars/<group>/all.yaml` (hierarchical group nodes)
- `group-vars/<group>/<child>/all.yaml` (deeper hierarchy nodes)
- `host_vars/<hostname>.yaml` (host-specific deltas and fallback conflicts)
- `hosts` (inventory mapping hosts to leaf groups)

Rules:

- Every group node is a directory containing `all.yaml`.
- Hosts belong to leaf groups only.
- Parent group inheritance is represented through hierarchy, not duplicate peer memberships.

## Merge Order (Runtime)

Runtime merge must be explicit and deterministic:

1. `group-vars/all.yaml`
2. Ancestor chain `all.yaml` files from root to leaf
3. Leaf group `all.yaml`
4. `host_vars/<hostname>.yaml`

After merge assembly, apply `_merge_directive` processing.

## Generation Algorithm

1. Build host payloads from parsed configurations.
2. Derive candidate common groups from host intersections (respecting min-host threshold).
3. Build a tree using strict host-set inclusion:
   - parent is the nearest strict superset.
4. Compute raw payload per group node.
5. Write each group node payload as strict delta against merged parent ancestry.
6. Write each host payload as strict delta against effective inherited group payload.

## Single-Owner Path Invariant

For any effective host path, ownership must be unambiguous.

- Shared equal values are hoisted to nearest common ancestor.
- Child nodes only contain incremental deltas.
- Peer nodes should not define conflicting values for the same path inherited by the same host.

## Sibling Conflict Handling

A sibling conflict means two peer groups define different values for the same path and at least one host would inherit both.

### Build-time policy

1. Attempt to resolve by restructuring ownership in this order:
   - Hoist common values to ancestor.
   - Push divergent values down into non-overlapping descendants.
2. If conflict still cannot be resolved by group restructuring, fallback is mandatory:
   - Remove the conflicting path from sibling group payloads for affected hosts.
   - Emit host-specific values into `host_vars/<hostname>.yaml` for the affected hosts.
3. If fallback cannot produce an unambiguous result, fail generation with diagnostics.

### Runtime policy

Runtime must reject any remaining sibling conflict.

- Encountering a sibling conflict during merge is an error.
- Do not use load-order precedence to silently pick one value.
- Error diagnostics must include:
  - conflicting path
  - sibling group names
  - affected host(s)

## Determinism Requirements

- Same inputs produce byte-stable YAML outputs.
- Group ordering and serialization are stable.
- Conflict detection and fallback decisions are deterministic.

## Validation Gates

Generation should enforce the following checks:

1. No unresolved runtime sibling conflicts remain.
2. Duplicate path ratio across group payloads is below configured threshold.
3. Host vars are minimal deltas (size trend tracked).
4. Determinism test passes across repeated runs.

## Non-Goals

- Supporting legacy `group_vars/*.yml` flat layouts.
- Preserving historical CLI behavior for deprecated hierarchy modes.
- Implementation-type specific grouping semantics.
