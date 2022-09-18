# hadolint ignore=DL3006
FROM gcr.io/distroless/static-debian11
COPY unisrv /
EXPOSE 5000
ENV UNISRV_HOST=0.0.0.0
ENTRYPOINT [ "/unisrv", "/app" ]
