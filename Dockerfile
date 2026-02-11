# Copyright 2024 The Kevin Berger <huhouhuam@gmail.com> Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http:www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.26.0-alpine AS build

LABEL maintainer="huhouhuam@gmail.com"

ENV GOPATH=/go
ENV CGO_ENABLED=0

RUN apk add -U --no-cache ca-certificates
RUN apk add -U curl
RUN curl -s -q https://raw.githubusercontent.com/huhouhua/glctl/main/LICENSE -o /go/LICENSE
RUN go install -v -ldflags "$(go run scripts/gen-ldflags.go)" "github.com/huhouhua/glctl@latest"

FROM alpine:3.18

COPY --from=build /go/bin/glctl  /usr/bin/glctl
COPY --from=build /go/LICENSE /licenses/LICENSE
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["glctl"]