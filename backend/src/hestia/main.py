from fastapi import FastAPI

app = FastAPI(title="Hestia API")


@app.get("/health")
async def health() -> dict[str, str]:
    return {"status": "ok"}
