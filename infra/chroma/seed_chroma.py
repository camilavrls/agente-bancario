import os
import pathlib
import time

import chromadb


COLLECTION_NAME = "banking_kb"
DEFAULT_DOCS_DIR = "/docs"
KEYWORDS = [
    "emprestimo",
    "consignado",
    "aposentado",
    "taxa",
    "tarifa",
    "conta",
    "ted",
    "pix",
    "limite",
    "cartao",
    "aumento",
    "risco",
    "confirmacao",
    "seguranca",
    "cliente",
    "autorizacao",
]


def normalize(text):
    return (
        text.lower()
        .replace("á", "a")
        .replace("à", "a")
        .replace("ã", "a")
        .replace("â", "a")
        .replace("é", "e")
        .replace("ê", "e")
        .replace("í", "i")
        .replace("ó", "o")
        .replace("õ", "o")
        .replace("ô", "o")
        .replace("ú", "u")
        .replace("ç", "c")
    )


def embed(text):
    normalized = normalize(text)
    vector = [float(normalized.count(keyword)) for keyword in KEYWORDS]

    norm = sum(value * value for value in vector) ** 0.5
    if norm == 0:
        return vector

    return [value / norm for value in vector]


def wait_for_chroma(client):
    for _ in range(30):
        try:
            client.heartbeat()
            return
        except Exception:
            time.sleep(1)

    raise RuntimeError("Chroma nao ficou disponivel a tempo")


def main():
    host = os.getenv("CHROMA_HOST", "localhost")
    port = int(os.getenv("CHROMA_PORT", "8000"))
    docs_dir = pathlib.Path(os.getenv("CHROMA_DOCS_DIR", DEFAULT_DOCS_DIR))

    client = chromadb.HttpClient(host=host, port=port)
    wait_for_chroma(client)

    collection = client.get_or_create_collection(
        name=COLLECTION_NAME,
        metadata={"description": "Base de conhecimento bancaria ficticia"},
    )

    print(f"Collection pronta: {COLLECTION_NAME}")

    documents = sorted(docs_dir.glob("*.md"))
    print(f"Documentos encontrados: {len(documents)}")
    for document in documents:
        print(f"- {document.name}")

    ids = []
    contents = []
    metadatas = []
    embeddings = []

    for document in documents:
        content = document.read_text(encoding="utf-8")
        ids.append(document.stem)
        contents.append(content)
        metadatas.append({"source": document.name})
        embeddings.append(embed(content))

    if ids:
        collection.upsert(
            ids=ids,
            documents=contents,
            metadatas=metadatas,
            embeddings=embeddings,
        )

    print(f"Documentos gravados no Chroma: {len(ids)}")

    smoke_test_query = "qual a taxa do emprestimo consignado?"
    smoke_test_result = collection.query(
        query_embeddings=[embed(smoke_test_query)],
        n_results=1,
        include=["documents", "metadatas"],
    )

    print(f"Smoke test query: {smoke_test_query}")
    print(f"Smoke test result: {smoke_test_result['metadatas'][0][0]['source']}")


if __name__ == "__main__":
    main()
