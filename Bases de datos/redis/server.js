const express = require('express');
const app = express();



const redisClient = require('./redis-client');
//client = redisClient.createClient();

app.get('/store/:key', async (req, res) => {
  const { key } = req.params;
  const value = req.query;
  console.log(key);
  console.log(value);
  // var multi = redisClient.multi();

  // multi.rpush(key, JSON.stringify(value));
  // await multi.exec(function(errors, results) {

  // })
    await redisClient.setAsync(key, JSON.stringify(value));
    //client.set(key, JSON.stringify(value));
  return res.send('Success');
});

app.get('/:key', async (req, res) => {
  const { key } = req.params;
  const rawData = await redisClient.getAsync(key);
  return res.json(rawData);
});

app.get('/', (req, res) => {
  return res.send('Hello world');
});

const PORT = process.env.PORT || 3500;
app.listen(PORT, () => {
  console.log(`Server listening on port ${PORT}`);
});


