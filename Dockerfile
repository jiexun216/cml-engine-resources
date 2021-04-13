FROM alpine:latest

ADD cml-engine-resources /cml-engine-resources
ENTRYPOINT ["./cml-engine-resources"]