FROM webhippie/alpine:latest
MAINTAINER Thomas Boerger <thomas@webhippie.de>

EXPOSE 8080
VOLUME ["/var/lib/ldap-proxy"]

RUN apk update && \
  apk add \
    ca-certificates \
    bash && \
  rm -rf \
    /var/cache/apk/* && \
  addgroup \
    -g 1000 \
    ldap-proxy && \
  adduser -D \
    -h /var/lib/ldap-proxy \
    -s /bin/bash \
    -G ldap-proxy \
    -u 1000 \
    ldap-proxy

COPY ldap-proxy /usr/bin/

USER ldap-proxy
ENTRYPOINT ["/usr/bin/ldap-proxy"]
CMD ["server"]

# ARG VERSION
# ARG BUILD_DATE
# ARG VCS_REF

# LABEL org.label-schema.version=$VERSION
# LABEL org.label-schema.build-date=$BUILD_DATE
# LABEL org.label-schema.vcs-ref=$VCS_REF
LABEL org.label-schema.vcs-url="https://github.com/webhippie/ldap-proxy.git"
LABEL org.label-schema.name="LDAP Proxy"
LABEL org.label-schema.vendor="Thomas Boerger"
LABEL org.label-schema.schema-version="1.0"
