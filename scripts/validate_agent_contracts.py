#!/usr/bin/env python3
"""Validate repo-local agent metadata against repository reality."""

from __future__ import annotations

import json
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
INDEX = ROOT / "agent/module-index.json"
GRAPH = ROOT / "agent/dependency-graph.json"
MODULES = ROOT / "agent/modules"
VALID_STATES = {"planned", "implemented", "deprecated"}
VALID_PROVENANCE = {"authored", "generated"}


def load(path: Path):
    return json.loads(path.read_text(encoding="utf-8"))


def main() -> int:
    errors: list[str] = []
    warnings: list[str] = []

    required = [
        ROOT / "AGENTS.md",
        ROOT / "agent/methodology.json",
        INDEX,
        GRAPH,
        ROOT / "agent/schemas/methodology.schema.json",
        ROOT / "agent/schemas/module-index.schema.json",
        ROOT / "agent/schemas/module-manifest.schema.json",
        ROOT / "agent/schemas/dependency-graph.schema.json",
    ]
    for path in required:
        if not path.exists():
            errors.append(f"missing required path: {path.relative_to(ROOT)}")

    if errors:
        return report(errors, warnings)

    try:
        methodology = load(ROOT / "agent/methodology.json")
        index = load(INDEX)
        graph = load(GRAPH)
    except (OSError, json.JSONDecodeError) as exc:
        errors.append(f"cannot load agent metadata: {exc}")
        return report(errors, warnings)

    if methodology.get("methodology") != "agent-project-methodology" or methodology.get("version") != "1.1":
        errors.append("agent/methodology.json must declare agent-project-methodology v1.1")

    modules = index.get("modules")
    if index.get("schema_version") != 1 or not isinstance(modules, dict):
        errors.append("invalid agent/module-index.json")
        return report(errors, warnings)

    known = set(modules)
    expected_manifests: set[Path] = set()

    for module_id, meta in modules.items():
        if meta.get("status") not in VALID_STATES:
            errors.append(f"{module_id}: invalid module status")
        manifest_text = meta.get("manifest")
        if not isinstance(manifest_text, str):
            errors.append(f"{module_id}: missing manifest path")
            continue
        manifest_path = ROOT / manifest_text
        expected_manifests.add(manifest_path.resolve())
        if not manifest_path.is_file():
            errors.append(f"{module_id}: missing manifest {manifest_text}")
            continue
        try:
            manifest = load(manifest_path)
        except (OSError, json.JSONDecodeError) as exc:
            errors.append(f"{module_id}: cannot load manifest: {exc}")
            continue
        if manifest.get("module_id") != module_id:
            errors.append(f"{module_id}: manifest module_id mismatch")
        paths = manifest.get("paths")
        if not isinstance(paths, list):
            errors.append(f"{module_id}: paths must be an array")
            continue
        for item in paths:
            if not isinstance(item, dict):
                errors.append(f"{module_id}: path entry must be an object")
                continue
            path_text = item.get("path")
            lifecycle = item.get("lifecycle")
            provenance = item.get("provenance")
            if not isinstance(path_text, str) or lifecycle not in VALID_STATES or provenance not in VALID_PROVENANCE:
                errors.append(f"{module_id}: invalid path metadata {item!r}")
                continue
            exists = (ROOT / path_text).exists()
            if lifecycle == "planned" and exists:
                errors.append(f"{module_id}: planned path exists: {path_text}")
            if lifecycle in {"implemented", "deprecated"} and not exists:
                errors.append(f"{module_id}: {lifecycle} path missing: {path_text}")
            if provenance == "generated" and (not item.get("generator") or not item.get("check")):
                errors.append(f"{module_id}: generated path lacks generator/check: {path_text}")
            if lifecycle == "deprecated" and not isinstance(item.get("legacy_consumers"), list):
                errors.append(f"{module_id}: deprecated path lacks legacy_consumers: {path_text}")

        if meta.get("status") == "implemented":
            for path_text in meta.get("roots", []):
                if not (ROOT / path_text).exists():
                    errors.append(f"{module_id}: implemented root missing: {path_text}")
            for path_text in meta.get("entrypoints", []):
                if not (ROOT / path_text).is_file():
                    errors.append(f"{module_id}: implemented entrypoint missing: {path_text}")

        for dep in meta.get("dependencies", []):
            if dep not in known:
                errors.append(f"{module_id}: unknown dependency {dep}")

    actual_manifests = {p.resolve() for p in MODULES.glob("*.json")}
    for orphan in sorted(actual_manifests - expected_manifests):
        errors.append(f"orphan manifest: {orphan.relative_to(ROOT)}")

    edges = graph.get("edges")
    if graph.get("schema_version") != 1 or not isinstance(edges, list):
        errors.append("invalid agent/dependency-graph.json")
        edges = []
    normalized: set[tuple[str, str]] = set()
    for edge in edges:
        if not isinstance(edge, list) or len(edge) != 2 or not all(isinstance(x, str) for x in edge):
            errors.append(f"invalid dependency edge: {edge!r}")
            continue
        source, target = edge
        if source not in known or target not in known:
            errors.append(f"unknown dependency edge: {source} -> {target}")
        normalized.add((source, target))

    for module_id, meta in modules.items():
        declared = set(meta.get("dependencies", []))
        graphed = {target for source, target in normalized if source == module_id}
        if declared != graphed:
            errors.append(f"{module_id}: index/graph dependencies differ: index={sorted(declared)} graph={sorted(graphed)}")

    if INDEX.stat().st_size > 12 * 1024:
        warnings.append("agent/module-index.json exceeds 12 KB warning budget")
    for path in MODULES.glob("*.json"):
        if path.stat().st_size > 10 * 1024:
            warnings.append(f"{path.relative_to(ROOT)} exceeds 10 KB warning budget")

    return report(errors, warnings)


def report(errors: list[str], warnings: list[str]) -> int:
    for warning in warnings:
        print(f"WARNING: {warning}")
    for error in errors:
        print(f"ERROR: {error}", file=sys.stderr)
    if errors:
        print(f"agent contracts: FAILED ({len(errors)} error(s))", file=sys.stderr)
        return 1
    print("agent contracts: OK")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
