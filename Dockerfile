FROM alpine:3.6

ADD ./k8s-kvm-health /k8s-kvm-health

ENTRYPOINT ["/k8s-kvm-health"]
CMD ["daemon"]
