FROM ubuntu:16.04
RUN apt-get update && apt-get install -y \
      curl \
      dbus \
      kmod \
      iproute2 \
      iputils-ping \
      net-tools \
      openssh-server \
      rng-tools \
      sudo \
      systemd \
      udev \
      vim-tiny \
      wget \
      nano && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
RUN echo "" > /etc/machine-id && echo "" > /var/lib/dbus/machine-id
RUN sed -i -e 's/^AcceptEnv LANG LC_\*$/#AcceptEnv LANG LC_*/' /etc/ssh/sshd_config
RUN sed -ri 's/^#?PermitRootLogin\s+.*/PermitRootLogin yes/' /etc/ssh/sshd_config
RUN echo "nameserver 8.8.8.8" > /etc/resolv.conf
RUN service ssh start
EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]