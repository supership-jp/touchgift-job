FROM fluent/fluentd:v1.13.3-debian-1.0

# Use root account to use apt
USER root

# below RUN includes plugin as examples elasticsearch is not required
# you may customize including plugins as you wish
RUN apt-get update \
  && apt-get install -y --no-install-recommends netcat binutils \
  && buildDeps="make gcc g++ libc-dev wget bzip2 gnupg dirmngr" \
  && apt-get install -y --no-install-recommends $buildDeps \
  && gem install fluent-plugin-prometheus \
  && gem install fluent-plugin-rewrite-tag-filter \
  && gem install resolv -v 0.2.1 \
  && gem sources --clear-all \
  && wget https://github.com/jemalloc/jemalloc/releases/download/3.6.0/jemalloc-3.6.0.tar.bz2 \
  && tar xvf jemalloc-3.6.0.tar.bz2 \
  && cd jemalloc-3.6.0 \
  && ./configure --prefix=/opt/jemalloc-3.6.0 \
  && make && make install && cd .. \
  && rm jemalloc-3.6.0.tar.bz2 && rm -rf jemalloc-3.6.0 \
  && apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false $buildDeps \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* \
  && rm -rf /tmp/* /var/tmp/* /usr/lib/ruby/gems/*/cache/*.gem

COPY config/* /fluentd/etc/

# NOTE: resolv v0.2.1 includes the fix for CPU spike issue due to DNS resolver.
# Forcing to load specific version of resolv (instead of bundled by default) is needed for Ruby 2.6.7, 2.7.3 and 3.0.1.
ENV RUBYLIB="/usr/local/bundle/gems/resolv-0.2.1/lib"
ENV LD_PRELOAD="/opt/jemalloc-3.6.0/lib/libjemalloc.so"
ENV RUBY_GC_HEAP_OLDOBJECT_LIMIT_FACTOR=0.9
ENV MALLOC_CONF="narenas:2"

# USER fluent
