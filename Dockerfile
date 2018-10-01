FROM larjim/kademlialab:latest
RUN mkdir -p /home/go/src/github.com/davgraols/D7024E-GR8/
RUN /usr/local/go/bin/go get -u github.com/golang/protobuf/proto
RUN /usr/local/go/bin/go get -u github.com/golang/protobuf/protoc-gen-go
COPY . /home/go/src/github.com/davgraols/D7024E-GR8/
WORKDIR /home/go/src/github.com/davgraols/D7024E-GR8/d7024e/
RUN /usr/local/go/bin/go build