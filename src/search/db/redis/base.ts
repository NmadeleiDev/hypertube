import * as redis from 'redis';
import log from '../../logger/logger';

export class BaseRedis {
  constructor() {
    this.initConnection();
  }

  private _client: redis.RedisClient | undefined;

  get client(): redis.RedisClient | undefined {
    return this._client;
  }

  set client(val: redis.RedisClient | undefined) {
    this._client = val;
  }

  get isConnected(): boolean {
    return !!this.client && this.client.connected;
  }

  private initConnection() {
    log.debug('process.env: ', process.env);

    const host =
      process.env.REDIS_HOST && process.env.REDIS_HOST.length > 0
        ? process.env.REDIS_HOST
        : 'localhost';
    const port =
      process.env.REDIS_PORT && process.env.REDIS_PORT.length > 0
        ? process.env.REDIS_PORT
        : '6379';
    const user =
      process.env.REDIS_USER && process.env.REDIS_USER.length > 0
        ? process.env.REDIS_USER
        : 'chat_user';
    const password =
      process.env.REDIS_PASSWORD && process.env.REDIS_PASSWORD.length > 0
        ? process.env.REDIS_PASSWORD
        : '234234234';

    const url = `redis://${user}:${password}@${host}:${port}`;
    log.info('Redis url: ', url);

    this.client = redis.createClient({
      url: url,
    });

    this.client
      .on('connect', function () {
        log.info('Connected to redis: ', this.isConnected);
      })
      .on('error', function (error) {
        log.error('Create client error: ', error);
      });
  }
}
