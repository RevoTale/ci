FROM ghcr.io/updatecli/updatecli:v0.123.0@sha256:5596883b1d8d8e174189a89467c9536b01f80ac26d260dc98eb2d19621caab31 AS updatecli

# Go keeps go.sum in sync; Git supports module downloads outside a Go proxy.
FROM golang:1.27.2-trixie@sha256:47667934874913a449b96faa7ef4777fa285d9ec29126ee80ed61c45b7cccf0e AS runtime
COPY --from=updatecli /usr/local/bin/updatecli /usr/local/bin/updatecli
ENTRYPOINT ["updatecli"]
CMD ["help"]

FROM runtime AS artifact-test
RUN ["go", "version"]
RUN ["git", "--version"]
RUN ["updatecli", "version"]

FROM runtime AS release
