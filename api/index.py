#!/usr/bin/env python
import requests
import six
import imagehash

from PIL import Image
from sanic import Sanic
from sanic.response import json

app = Sanic('CQT Memes')


def alpha_remover(image):
    if image.mode != 'RGBA':
        return image
    canvas = Image.new('RGBA', image.size, (255, 255, 255, 255))
    canvas.paste(image, mask=image)
    return canvas.convert('RGB')


def image_loader(content):
    image = alpha_remover(Image.open(six.BytesIO(content)))
    return imagehash.colorhash(image)


@app.route('/hash')
async def index(request):
    image_hash = ''
    url = request.args.get('url')
    if url:
        response = requests.get(url)
        if response.status_code == 200:
            image_hash = image_loader(response.content)
    return json({
      'url': url,
      'hash': str(image_hash)})
