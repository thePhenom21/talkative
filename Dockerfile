# Copyright 2022 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Use the offical Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.21 as builder
WORKDIR /app

# Initialize a new Go module.
COPY go.mod go.sum ./
RUN go mod download


# Copy local code to the container image.
COPY *.go ./
COPY .env ./
COPY weaver.toml ./

EXPOSE 8082/tcp

RUN go install
RUN go install github.com/ServiceWeaver/weaver/cmd/weaver@latest

# Build the command inside the container.
RUN CGO_ENABLED=0 GOOS=linux go build -o /app

ENV PATH="/go/bin:$PATH"


CMD ["weaver","single","deploy","weaver.toml"]
