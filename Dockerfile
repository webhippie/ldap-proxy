FROM webhippie/alpine:latest

LABEL maintainer="Thomas Boerger <thomas@webhippie.de>" \
  org.label-schema.name="LDAP Proxy" \
  org.label-schema.vendor="Thomas Boerger" \
  org.label-schema.schema-version="1.0"

EXPOSE 80 443 9000
VOLUME ["/var/lib/ldap-proxy"]

ENV LDAP_PROXY_HEALTH_ADDR 0.0.0.0:9000
ENV LDAP_PROXY_SERVER_STORAGE /var/lib/ldap-proxy
ENV LDAP_PROXY_SERVER_TEMPLATES /usr/share/ldap-proxy/templates
ENV LDAP_PROXY_SERVER_ASSETS /usr/share/ldap-proxy/assets

ENTRYPOINT ["/usr/bin/ldap-proxy"]
CMD ["server"]

RUN apk add --no-cache ca-certificates mailcap bash

COPY assets /usr/share/ldap-proxy/
COPY templates /usr/share/ldap-proxy/
COPY dist/binaries/ldap-proxy-*-linux-amd64 /usr/bin/ldap-proxy
