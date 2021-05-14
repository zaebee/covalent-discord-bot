#!/usr/bin/env python
import json
import requests
import six
import imagehash

from PIL import Image
from http.server import BaseHTTPRequestHandler


def alpha_remover(image):
    if image.mode != 'RGBA':
        return image
    canvas = Image.new('RGBA', image.size, (255, 255, 255, 255))
    canvas.paste(image, mask=image)
    return canvas.convert('RGB')


def image_loader(content):
    image = alpha_remover(Image.open(six.BytesIO(content)))
    return imagehash.colorhash(image)


class handler(BaseHTTPRequestHandler):

  def do_GET(self):
    self.send_response(200)
    self.send_header('Content-type', 'text/plain')
    self.end_headers()
    self.wfile.write('test')
    return


def hash_image(url):
    image_hash = ''
    # url = request.args.get('url')
    if url:
        response = requests.get(url)
        if response.status_code == 200:
            image_hash = image_loader(response.content)
    return json.dumps({
      'url': url,
      'hash': str(image_hash)})
