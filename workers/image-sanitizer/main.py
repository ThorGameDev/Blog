import os

from fastapi import FastAPI
from fastapi.responses import PlainTextResponse
from PIL import Image

app = FastAPI()

@app.get("/request", response_class=PlainTextResponse)
def request(file: str):
    # get the file path of the image to work on
    fullPath = os.path.abspath(os.path.join("./workdir/", file))

    # Ensure the output directory exists
    if not os.path.exists("./workdir/out/"):
        os.makedirs("./workdir/out/")

    # Create an output file path with a .jpeg extension
    baseFileName, _ = os.path.splitext(file)
    outFile = baseFileName+".jpeg"
    outpath = os.path.abspath(os.path.join("./workdir/out/", outFile))

    # Open the file as an image
    img = Image.open(fullPath)

    # Convert to RGB
    img = img.convert("RGB")

    # Minimize the image
    img.thumbnail((256, 256), Image.LANCZOS)

    # Save as minimal JEPG
    img.save(outpath, format="JPEG", quality=85, optimize=True)

    # Remove original image
    os.remove(fullPath)

    # Inform the backend the new file name
    return outFile
