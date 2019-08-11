#!/bin/sh
gunicorn -b 0.0.0.0:8000 -k uvicorn.workers.UvicornWorker example:app