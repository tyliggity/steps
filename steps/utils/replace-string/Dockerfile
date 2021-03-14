FROM python:3

WORKDIR /src/

ARG STEP_BASEPATH

COPY ${STEP_BASEPATH} /src/

RUN pip install black
RUN black --check .
RUN pip install -r ./requirements.txt

ENTRYPOINT ["./replace_string.py"]
