const MONGO_USER = process.env.MONGO_USER || 'admin';
const MONGO_PASSWD = process.env.MONGO_PASSWORD || 'passwd';
const MONGO_ADDR = process.env.MONGO_ADDRESS || 'localhost:27017';

export const DSN = `mongodb://${MONGO_USER}:${MONGO_PASSWD}@${MONGO_ADDR}/`;
export const usersDb = `hypertube`;
export const usersCollection = `search`;

console.log(DSN);
