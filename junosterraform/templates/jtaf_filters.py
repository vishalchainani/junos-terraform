#!/usr/bin/env python3
"""JTAF Ansible filters for merge-directive aware payload handling."""

from __future__ import annotations

from copy import deepcopy
from typing import Any


def _deep_merge(base: Any, override: Any) -> Any:
    """Recursively merge dict/list values, preferring override leaves."""
    if isinstance(base, dict) and isinstance(override, dict):
        merged = deepcopy(base)
        for key, value in override.items():
            if key in merged:
                merged[key] = _deep_merge(merged[key], value)
            else:
                merged[key] = deepcopy(value)
        return merged

    if isinstance(base, list) and isinstance(override, list):
        # For list-of-dict payloads keyed by "name", merge by key to avoid
        # duplicating nested structures like interface.unit and syslog.contents.
        if all(isinstance(item, dict) and "name" in item for item in base + override):
            return _merge_list_by_key(deepcopy(base) + deepcopy(override), "name")
        return deepcopy(base) + deepcopy(override)

    return deepcopy(override)


def _merge_list_by_key(values: list[Any], merge_key: str) -> list[Any]:
    """Merge list entries of dicts by a stable key while preserving order."""
    merged_by_key: dict[str, Any] = {}
    order: list[str] = []
    passthrough: list[Any] = []

    for item in values:
        if not isinstance(item, dict) or merge_key not in item:
            passthrough.append(deepcopy(item))
            continue

        key_value = str(item[merge_key])
        if key_value not in merged_by_key:
            merged_by_key[key_value] = deepcopy(item)
            order.append(key_value)
        else:
            merged_by_key[key_value] = _deep_merge(merged_by_key[key_value], item)

    return [merged_by_key[k] for k in order] + passthrough


def _apply_list_directives(data: Any) -> Any:
    """Apply _merge_list_directives recursively to dict/list structures."""
    if isinstance(data, list):
        return [_apply_list_directives(item) for item in data]

    if not isinstance(data, dict):
        return data

    list_directives = data.get("_merge_list_directives", {})
    result: dict[str, Any] = {}

    for key, value in data.items():
        if key == "_merge_list_directives":
            result[key] = value
            continue

        processed_value = _apply_list_directives(value)
        directive = list_directives.get(key)

        if (
            isinstance(directive, dict)
            and directive.get("_merge_directive") == "merge_by_key"
            and isinstance(processed_value, list)
        ):
            merge_key = directive.get("_merge_key", "name")
            processed_value = _merge_list_by_key(processed_value, merge_key)

        result[key] = processed_value

    return result


class FilterModule:
    """Expose JTAF custom filters to Ansible."""

    def filters(self) -> dict[str, Any]:
        return {
            "jtaf_extract_directive": self.extract_directive,
            "jtaf_remove_meta": self.remove_meta_keys,
            "jtaf_apply_merge_directives": self.apply_merge_directives,
        }

    @staticmethod
    def extract_directive(data: Any) -> str:
        if isinstance(data, dict):
            return str(data.get("_merge_directive", "replace"))
        return "replace"

    @staticmethod
    def remove_meta_keys(data: Any) -> Any:
        """Strip merge metadata and transient markers from payloads."""
        if isinstance(data, dict):
            return {
                key: FilterModule.remove_meta_keys(value)
                for key, value in data.items()
                if not key.startswith("_merge") and key != "_applied_directive"
            }

        if isinstance(data, list):
            return [FilterModule.remove_meta_keys(item) for item in data]

        return data

    def apply_merge_directives(self, jtaf_effective: dict[str, Any]) -> dict[str, Any]:
        """Apply supported merge/list directives within an already merged payload."""
        processed = _apply_list_directives(deepcopy(jtaf_effective))
        return self._mark_directives(processed)

    def _mark_directives(self, data: Any) -> Any:
        """Mark nodes that carried explicit _merge_directive metadata."""
        if isinstance(data, dict):
            directive = data.get("_merge_directive")
            marked = {
                key: self._mark_directives(value)
                for key, value in data.items()
            }
            if directive:
                marked["_applied_directive"] = directive
            return marked

        if isinstance(data, list):
            return [self._mark_directives(item) for item in data]

        return data
