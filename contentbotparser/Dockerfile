FROM python:3.10

WORKDIR /pars

COPY requirements.txt ./

RUN pip install -r requirements.txt

COPY . ./

EXPOSE 8001