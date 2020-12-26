FROM centos:8
#COPY resolv.conf /etc/
RUN yum -y install \
        iproute \
        iputils \
        openssh-server \
        net-tools \
        procps-ng \
        sudo \
        wget && \
    yum clean all
RUN echo "nameserver 8.8.8.8" > /etc/resolv.conf
RUN ssh-keygen -q -t ecdsa -f /etc/ssh/ssh_host_ecdsa_key -C '' -N ''
RUN ssh-keygen -q -t rsa -f /etc/ssh/ssh_host_rsa_key -C '' -N ''
RUN ssh-keygen -q -t ed25519 -f /etc/ssh/ssh_host_ed25519_key -C '' -N ''
RUN echo "root:root" | chpasswd
EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]