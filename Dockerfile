FROM larjim/kademlialab:latest
RUN mkdir -p /home/go/src/github.com/davgraols/D7024E-GR8/
#RUN ["/usr/local/go/bin/go", "get", "-u", "github.com/golang/protobuf/{proto,protoc-gen-go}"]
RUN /usr/local/go/bin/go get -u github.com/golang/protobuf/proto
RUN /usr/local/go/bin/go get -u github.com/golang/protobuf/protoc-gen-go
COPY . /home/go/src/github.com/davgraols/D7024E-GR8/
WORKDIR /home/go/src/github.com/davgraols/D7024E-GR8/d7024e/
RUN /usr/local/go/bin/go build
#RUN chmod 755 bootstrap.sh
#ENTRYPOINT [ "./bootstrap.sh" ]
#ENTRYPOINT ["/usr/local/go/bin/go", "run", "main.go"]