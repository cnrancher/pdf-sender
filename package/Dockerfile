FROM alpine:3
COPY bin/pdf-sender /pdf-sender
COPY package/init.yml /etc/pdf-sender.yml
ENTRYPOINT ["/pdf-sender"]
CMD [ "-f", "/etc/pdf-sender.yml", "run" ]