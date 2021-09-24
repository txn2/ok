# OK

A Simple web server that returns a bit of useful JSON data for testing and
diagnosing ingress with Kubernetes.


### Release
```bash
goreleaser --skip-publish --rm-dist --skip-validate
```

```bash
GITHUB_TOKEN=$GITHUB_TOKEN goreleaser --rm-dist
```