# kaniko-revanced Release Process

This document outlines the release process and cadence for **kaniko-revanced**.

## Release Schedule & Cadence

1. **Bi-Weekly Regular Releases**:
   - Every two weeks, an automated release is triggered incorporating dependency upgrades, bug fixes, performance optimizations, and security patches.
   - Tagged as `v1.X.Y` semver format.

2. **On-Demand CVE & Security Hotfixes**:
   - For Critical or High severity CVEs in direct dependencies or base images, patch releases (e.g. `v1.X.Y+1`) are triggered and published immediately upon fix verification.

---

## Published Artifacts & Registry

Container images are published to Docker Hub repository **`kanikorevanced/executor`** and **`kanikorevanced/warmer`**:
- `kanikorevanced/executor:latest`, `:0.1.0` (Standard Executor)
- `kanikorevanced/executor:debug`, `:debug-0.1.0` (Debug with Busybox shell)
- `kanikorevanced/executor:slim`, `:slim-0.1.0` (Minimal footprint)
- `kanikorevanced/warmer:latest`, `:0.1.0` (Cache Warmer)

Every release includes:
- Multi-architecture container images (`linux/amd64`, `linux/arm64`, `linux/s390x`, `linux/ppc64le`)
- Cosign keyless signatures & attestation
- Syft Software Bill of Materials (SBOM) attached to each image
- Standalone static binaries attached as GitHub Release assets

---

## How to Trigger a Release

### Automated Workflow Trigger (Recommended)
1. Ensure all tests in the `main` branch pass.
2. In GitHub Actions, navigate to **Release kaniko-revanced** workflow.
3. Click **Run workflow**, specify the new version tag (e.g., `v0.1.0`), and execute.

### Tag-Based Release
Pushing a git tag formatted as `v[0-9]+.[0-9]+.[0-9]+*` will automatically trigger the release pipeline:
```bash
git checkout main
git pull origin main
git tag v0.1.0
git push origin v0.1.0
```
The GitHub Action will:
1. Compile multi-arch binaries.
2. Build multi-arch container images.
3. Push images to Docker Hub.
4. Sign images with Cosign.
5. Generate and attach Syft SBOMs.
6. Create the GitHub Release with changelog notes.

