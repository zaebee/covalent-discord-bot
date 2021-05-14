#!/usr/bin/env python
import requests
import six
import numpy as np
import imagehash

from PIL import Image
from sanic import Sanic
from sanic.response import json

app = Sanic()


def alpharemover(image):
    if image.mode != 'RGBA':
        return image
    canvas = Image.new('RGBA', image.size, (255, 255, 255, 255))
    canvas.paste(image, mask=image)
    return canvas.convert('RGB')


def image_loader(hashfunc):
    def function(content):
        image = alpharemover(Image.open(six.BytesIO(content)))
        return hashfunc(image)
    return function


def with_ztransform_preprocess(hashfunc, hash_size=8):
    def function(content):
        image = alpharemover(Image.open(six.BytesIO(content)))
        image = image.convert("L").resize((hash_size, hash_size), Image.ANTIALIAS)
        data = image.getdata()
        quantiles = np.arange(100)
        quantiles_values = np.percentile(data, quantiles)
        zdata = (np.interp(data, quantiles_values, quantiles) / 100 * 255).astype(np.uint8)
        image.putdata(zdata)
        return hashfunc(image)
    return function


def hash_from_url(url):
    response = requests.get(url)
    if response.status_code == 200:
        return image_loader(imagehash.colorhash)(response.content)


@app.route('/')
async def index(request):
    image_hash = ''
    url = request.args.get('url')
    if url:
        image_hash = hash_from_url(url)
    return json({
      'url': url,
      'hash': str(image_hash)})
