import { Pool, QueryConfig, QueryResult } from 'pg';
import { DSN } from './config';
import log from '../../logger/logger';

export class PostgresUser {
  constructor() {
    this.initConnection().catch((e) => log.warn(e));
  }

  private static _pool = new Pool({
    max: 20,
    connectionString: DSN,
    idleTimeoutMillis: 30000,
  });

  static async query(text: string = 'SELECT NOW()', values = []) {
    const client = await PostgresUser._pool.connect();
    log.debug('got client: ', client);
    try {
      const res = await client.query(text, values);
      console.log('PostgresUser result: ', res);
      return res;
    } finally {
      // Make sure to release the client before any error handling,
      // just in case the error handling itself throws an error.
      client.release();
    }
  }

  static async getClient() {
    return await PostgresUser._pool.connect();
  }

  async initConnection() {
    try {
      PostgresUser._pool = new Pool({
        max: 20,
        connectionString: DSN,
        idleTimeoutMillis: 30000,
      });
      log.info('Postrgres connection succeded!');
    } catch (e) {
      log.error('Failed to connect to Postrgres: ', e);
      throw 'Connection error!';
    }
  }
}
