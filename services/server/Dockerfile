FROM python:3.7.3-slim

ENV FRONTEND_ORIGIN "http://localhost:8080/"
# install gcc
RUN apt-get update && \
    apt-get -y install gcc && \
    apt-get clean

# set working directory
WORKDIR /app
ENV PYTHONPATH=/app
# add and install requirements
COPY ./requirements.txt /app/requirements.txt
RUN pip install -r requirements.txt

# add entrypoint.sh
COPY ./entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# add gunicorn config
COPY ./gunicorn_conf.py /gunicorn_conf.py

# add app
COPY . /app

# run server
CMD ["/app/entrypoint.sh"]