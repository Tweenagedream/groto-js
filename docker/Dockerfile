FROM alpine:edge
RUN apk update && apk add ca-certificates && update-ca-certificates && apk add openssl
RUN apk add openjdk7-jre
RUN apk add protobuf
WORKDIR /root
RUN wget https://github.com/google/protobuf/releases/download/v3.1.0/protobuf-js-3.1.0.tar.gz
RUN wget https://dl.google.com/closure-compiler/compiler-latest.zip
RUN wget https://github.com/google/closure-library/archive/v20161024.zip

RUN for each in $(ls *.zip); do unzip $each; done
RUN tar xzf *.tar.gz

RUN mv closure-library* closure-library
ADD entrypoint.sh /root/entrypoint.sh
RUN chmod +x /root/entrypoint.sh
ENTRYPOINT ['/root/entrypoint.sh']
