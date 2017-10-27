FROM webhippie/alpine:latest
MAINTAINER Thomas Boerger <thomas@webhippie.de>

EXPOSE 8080 80 443
VOLUME ["/var/lib/ldap-proxy"]

LABEL org.label-schema.version=latest
LABEL org.label-schema.name="LDAP Proxy"
LABEL org.label-schema.vendor="Thomas Boerger"
LABEL org.label-schema.schema-version="1.0"

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
    ldap-proxy && \
  mkdir -p \
    /usr/share/ldap-proxy

ENV LDAP_PROXY_SERVER_STORAGE /var/lib/ldap-proxy
ENV LDAP_PROXY_SERVER_TEMPLATES /usr/share/ldap-proxy/templates
ENV LDAP_PROXY_SERVER_ASSETS /usr/share/ldap-proxy/assets

COPY assets /usr/share/ldap-proxy/
COPY templates /usr/share/ldap-proxy/
COPY ldap-proxy /usr/bin/

USER ldap-proxy
ENTRYPOINT ["/usr/bin/ldap-proxy"]
CMD ["server"]
