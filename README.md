# Memes bot

A bot that is listening discord channel and stores messages with attachments (memes) on frontend website.

## Build Setup

Copy and update settings in `.env.example`
```shell
  set DISCORD_TOKEN=<bot_token>
  set DISCORD_CHANNEL=<channel_id>
  set ELASTICSEARCH_URL=localhost:9200
  set FRONTEND_URL=https://covalent.me-mes.site
  set IMAGE_HASH_URL=localhost:8000
```

### Frontend

```bash
# install dependencies
$ npm install

# serve with hot reload at localhost:3000
$ npm run dev

# build for production and launch server
$ npm run build
$ npm run start

# generate static project
$ npm run generate
```

For detailed explanation on how things work, check out [Nuxt.js docs](https://nuxtjs.org).

### Backend installation

After that you can run discord bot to listen messages.

```bash
# install dependencies for python script that checks imageHash. 
$ pip install -r requirements.txt

# run imageHash service.
$ cd api; sanic index.app -H 127.0.0.1 -p 8000

# run discord bot
$ cd bot; go run main.go
```
