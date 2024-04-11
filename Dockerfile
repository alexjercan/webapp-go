FROM postgres

RUN apt-get update && apt-get install -y postgresql-contrib
