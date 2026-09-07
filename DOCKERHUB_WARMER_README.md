# kaniko-revanced: Cache Warmer

[![GitHub Stars](https://img.shields.io/github/stars/raj-sh-git/kaniko-revanced?style=social)](https://github.com/raj-sh-git/kaniko-revanced)
[![Docker Pulls](https://img.shields.io/docker/pulls/kanikorevanced/warmer.svg)](https://hub.docker.com/r/kanikorevanced/warmer)

**kaniko-revanced warmer** (`kanikorevanced/warmer`) pre-fetches and populates base container image layers into a local volume cache or remote cache repository to accelerate subsequent Kaniko builds.

---

## 🏷️ Available Tags

* **`0.2.0`**, **`latest`**: Multi-arch (`linux/amd64`, `linux/arm64`, `linux/s390x`, `linux/ppc64le`) cache warmer.

---

## 🛠️ Usage Example

### Pre-populating a Cache Volume in Kubernetes

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: kaniko-cache-warmer
spec:
  containers:
  - name: warmer
    image: kanikorevanced/warmer:latest
    args:
      - "--cache-dir=/cache"
      - "--image=golang:1.24"
      - "--image=node:20-alpine"
      - "--image=alpine:3.20"
    volumeMounts:
      - name: kaniko-cache
        mountPath: /cache
  restartPolicy: Never
  volumes:
    - name: kaniko-cache
      persistentVolumeClaim:
        claimName: kaniko-cache-pvc
```

---

## 📄 License & Source Code

* **GitHub Repository**: [https://github.com/raj-sh-git/kaniko-revanced](https://github.com/raj-sh-git/kaniko-revanced)
* **License**: Apache License 2.0
