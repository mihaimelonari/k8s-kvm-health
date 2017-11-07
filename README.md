[![CircleCI](https://circleci.com/gh/giantswarm/k8s-kvm-health/tree/master.svg?style=shield)](https://circleci.com/gh/giantswarm/k8s-kvm-health/tree/master)
# k8s-kvm-health

Flannel-network-health serves as health endpoint for network configuration created by [flannel-operator](https://github.com/giantswarm/flannel-operator).

* endpoint `/bridge-healthz` - check if interface `br-${CLUSTER_ID}` is present and if it has configured right ip address
* endpoint `/flannel-healthz`- check if interface `flannel.${VNI}` is present and if it has configured right ip address


### How to build

#### Dependencies

Dependencies are managed using [`glide`](https://github.com/Masterminds/glide) and contained in the `vendor` directory. See `glide.yaml` for a list of libraries this project directly depends on and `glide.lock` for complete information on all external libraries and their versions used.

**Note:** The `vendor` directory is **flattened**. Always use the `--strip-vendor` (or `-v`) flag when working with `glide`.

#### Building the standard way

```nohighlight
go build
```

#### Cross-compiling in a container

Here goes the documentation on compiling for different architectures from inside a Docker container.

## Running PROJECT

- How to use
- What does it do exactly

## Contact

- Mailing list: [giantswarm](https://groups.google.com/forum/!forum/giantswarm)
- IRC: #[giantswarm](irc://irc.freenode.org:6667/#giantswarm) on freenode.org
- Bugs: [issues](https://github.com/giantswarm/PROJECT/issues)

## License

PROJECT is under the Apache 2.0 license. See the [LICENSE](/giantswarm/example-opensource-repo/blob/master/LICENSE) file for details.
