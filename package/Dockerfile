FROM registry.suse.com/bci/bci-base:15.3
RUN zypper -n install timezone &&\
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone
COPY bin/pdf-sender /pdf-sender
COPY package/init.yml /etc/pdf-sender.yml
ENTRYPOINT ["/pdf-sender"]
CMD [ "-f", "/etc/pdf-sender.yml", "run" ]