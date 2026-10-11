from sqlalchemy import text

from hestia.db import engine

async def test_postgis_is_available():
    async with engine.connect() as conn:
        result = await conn.execute(text("SELECT postgis_version()"))
        version = result.scalar_one()

    assert version is not None
