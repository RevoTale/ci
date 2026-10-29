FROM ghcr.io/updatecli/updatecli:v0.123.0 AS updatecli

# Go keeps go.sum in sync; Git supports module downloads outside a Go proxy.
FROM golang:1.27.2-trixie AS runtime
COPY --from=updatecli /usr/local/bin/updatecli /usr/local/bin/updatecli
ENTRYPOINT ["updatecli"]
CMD ["help"]

FROM runtime AS artifact-test
RUN ["go", "version"]
RUN ["git", "--version"]
RUN ["updatecli", "version"]

FROM runtime AS release
