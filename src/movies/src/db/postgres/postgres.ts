import { Pool, QueryResult } from 'pg';
import { DSN } from './config';
import log from '../../logger/logger';
const POSTGRES_SCHEME = process.env.POSTGRES_SCHEME || 'hypertube';

log.debug('creating pool');
const pool = new Pool({
  max: 20,
  connectionString: DSN,
  idleTimeoutMillis: 30000,
});
log.debug('pool created', pool);
log.debug('testing connection');
pool.query(`SELECT 1 from ${POSTGRES_SCHEME}.movies`, (error) => {
  if (error) throw error;
  console.log('Connection succeded');
});

export async function query(
  text: string,
  values: (string | number | Date)[] = []
): Promise<QueryResult<any>> {
  try {
    const start = Date.now();
    log.trace('Executing query', { text, values });
    const res = await pool.query({ text, values });
    const duration = Date.now() - start;
    log.info('Executed query', {
      text,
      values,
      duration,
      rows: res.rowCount,
    });
    log.trace('Result:', res.rows);
    return res;
  } catch (e) {
    log.error(e, text, values);
    return {
      rows: [],
      command: text.split(' ')[0],
      rowCount: 0,
      oid: 0,
      fields: [],
    };
  }
}

export const getClient = async () => {
  const client = await pool.connect();
  const query = client.query;
  const release = client.release;
  let lastQuery = null;
  // set a timeout of 5 seconds, after which we will log this client's last query
  const timeout = setTimeout(() => {
    log.error('A client has been checked out for more than 5 seconds!');
    log.error(`The last executed query on this client was: ${lastQuery}`);
  }, 5000);
  // monkey patch the query method to keep track of the last query executed
  client.query = (...args) => {
    lastQuery = args;
    return query.apply(client, args);
  };
  client.release = () => {
    // clear our timeout
    clearTimeout(timeout);
    // set the methods back to their old un-monkey-patched version
    client.query = query;
    client.release = release;
    return release.apply(client);
  };
  return client;
};
