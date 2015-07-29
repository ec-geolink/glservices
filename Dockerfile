# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
# docker run -d -p 6789:6789 -p 8081:8081 [IMAGE ID] 
FROM golang

# Copy the local package files to the container's workspace.
RUN mkdir /go/src/geolink.org/

# Create a non-root user to run as
RUN groupadd -r gorunner -g 433 && \
mkdir /home/gorunner && \
useradd -u 431 -r -g gorunner -d /home/gorunner -s /sbin/nologin -c "User to run go apps on high ports" gorunner && \ 
chown -R gorunner:gorunner  /home/gorunner  && \ 
chown -R gorunner:gorunner /go/src/geolink.org

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get github.com/emicklei/go-restful
RUN go get github.com/linkeddata/gojsonld
RUN go get github.com/parnurzeal/gorequest
RUN go get github.com/knakk/sparql
# RUN go install data.oceandrilling.org/codices

# Need to replace this with a build and run
# set user
USER gorunner

# Move to a workign directory for running codices so it can see it's static files
# future version should take this as a param so static content can be anywhere
WORKDIR /go/src/geolink.org

RUN  git clone https://github.com/ec-geolink/glservices.git

WORKDIR /go/src/geolink.org/glservices

# Run the command by default when the container starts.
ENTRYPOINT go run main.go

# Document that the service listens on this port
EXPOSE 6789 8081