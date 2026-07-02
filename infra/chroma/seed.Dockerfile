FROM python:3.12-slim

RUN pip install --no-cache-dir chromadb

WORKDIR /app

COPY seed_chroma.py /app/seed_chroma.py

ENTRYPOINT ["python", "/app/seed_chroma.py"]
