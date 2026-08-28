# Security Policy for kaniko-revanced

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| latest  | :white_check_mark: |
| < 1.24  | :x:                |

## Reporting a Vulnerability

We take the security of **kaniko-revanced** seriously. Because kaniko runs inside Kubernetes and CI/CD environments with high privileges (often root userspace execution), rapid remediation of CVEs in dependencies and container base images is our top priority.

To report a vulnerability:
1. **GitHub Security Advisory (Recommended)**: Use the "Report a vulnerability" button under the **Security** tab of the `raj-sh-git/kaniko-revanced` GitHub repository to submit a private report.
2. Please include:
   - Detailed description of the vulnerability / CVE ID
   - Minimal reproduction Dockerfile / setup
   - Impact assessment and suggested remediation if known

### Response & Patching SLA
- **Acknowledgment**: Within 24-48 hours.
- **Critical / High CVEs**: Expedited patch release built and published to Docker Hub within 72 hours.
- **Standard CVEs & Dep Upgrades**: Bundled into our regular bi-weekly release cycle.

