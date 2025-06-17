import asyncio
import os

from fastapi import FastAPI, Query, Response
from fastapi.responses import JSONResponse
from generator import generate_attachment
from concurrent.futures import ProcessPoolExecutor


app = FastAPI()


@app.get("/gen")
async def generate(pwd: str = Query(...), id: str = Query(...), flags: str = Query(...)):
    if pwd != os.getenv("generator_pwd"):
        return JSONResponse(content={'error': 'Invalid password'}, status_code=403)

    try:
        loop = asyncio.get_running_loop()
        file = await loop.run_in_executor(ProcessPoolExecutor(), generate_attachment, flags.encode())
    except Exception as e:
        return JSONResponse(content={'error': str(e)}, status_code=500)

    return Response(content=file, media_type="application/zip", headers={
        'Content-Disposition': f'attachment; filename={id}.zip'
    })
