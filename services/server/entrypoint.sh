#!/bin/sh
gunicorn -b 0.0.0.0:8000 --worker-class=gevent --worker-connections=1000 example:app