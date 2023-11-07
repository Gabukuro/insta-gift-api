FROM golang as builder

WORKDIR /usr/application
COPY . .

ENV GO111MODULE=on
ENV CGO_ENABLED=0

RUN make install
RUN make build-prod

FROM scratch

ARG ENV
ENV ENV=$ENV

COPY --from=builder /tmp/main /main

EXPOSE 3000
CMD ["/main"]
