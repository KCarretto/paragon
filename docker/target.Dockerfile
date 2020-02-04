FROM ubuntu:16.04 as dev
WORKDIR /app
RUN apt-get update && \
    apt install -y openssh-server && \
    echo 'root:changeme' | chpasswd && \
    sed -i 's/PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config && \
    sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd && \
    mkdir -p /var/run/sshd && \
    chmod 0755 /var/run/sshd && \
    chown root:root /var/run/sshd

EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]
