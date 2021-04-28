import mongo from 'mongodb';
import log from '../../logger/logger';
import { DSN } from './config';

export class MongoUser {
  constructor() {
    this.initConnection().catch((e) => log.warn(e));
  }

  private _connection: mongo.MongoClient | undefined;

  get connection(): mongo.MongoClient | undefined {
    return this._connection;
  }

  set connection(val: mongo.MongoClient | undefined) {
    this._connection = val;
  }

  async initConnection() {
    try {
      const mongoClient = new mongo.MongoClient(DSN, {
        useUnifiedTopology: true,
      });
      this._connection = await mongoClient.connect();
      log.info('Mongo connection succeded!');
    } catch (e) {
      log.error('Failed to connect to mongo: ', e);
      throw 'Connection error!';
    }
  }
}

// export const MongoManager = new MongoUser()
