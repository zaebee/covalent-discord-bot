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


@app.route('/hash')
async def index(request):
    image_hash = ''
    url = request.args.get('url')
    if url:
        response = requests.get(url)
        if response.status_code == 200:
            image = alpha_remover(Image.open(six.BytesIO(response.content)))
            image_hash = imagehash.colorhash(image)
    return json({
        'url': url,
        'hash': str(image_hash)})
