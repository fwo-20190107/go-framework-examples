FROM mysql:8.0-debian

RUN apt-get update \
    && apt-get install locales -y \
    && sed -i -E 's/# (ja_JP.UTF-8)/\1/' /etc/locale.gen \
    && locale-gen \
    && update-locale LANG=ja_JP.UTF-8 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

ENV LC_ALL ja_JP.UTF-8