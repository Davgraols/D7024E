FROM larjim/kademlialab:latest
RUN mkdir -p /home/go/src/github.com/davgraols/D7024E-GR8
COPY . /home/go/src/github.com/davgraols/D7024E-GR8
WORKDIR /home/go/src/github.com/davgraols/D7024E-GR8
ENTRYPOINT ["/usr/local/go/bin/go", "run", "main.go"]