FROM python:3.14.7-alpine3.24

WORKDIR /home/worker/app

RUN pip install --upgrade pip

RUN addgroup -S worker --gid 1000 && adduser -S worker -G worker -u 1000

EXPOSE 8100

USER worker

COPY --chown=worker:worker requirements.txt .

RUN --mount=type=cache,target=/home/worker/.cache/pip,uid=1000,gid=1000 \
    pip install -r requirements.txt

ENV PATH="/home/worker/.local/bin:${PATH}"

COPY --chown=worker:worker ./main.py .

CMD ["fastapi", "run", "main.py",  "--port", "8100"]
