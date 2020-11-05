FROM scratch
#FROM centos:7

COPY config-example.yaml  /config.yaml
COPY helm-wrapper /helm-wrapper

EXPOSE 8080

CMD [ "/helm-wrapper" ]
