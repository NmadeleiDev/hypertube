import { httpResponse } from '../model/model';
import log from '../logger/logger';
const http = require('http');
const https = require('https');
const { readFile, createWriteStream, unlink } = require('fs');
import { writeFile } from 'fs/promises';

export function download(url: string): Promise<Buffer> {
  log.info('[Download]', url);
  return new Promise((resolve, reject) => {
    try {
      const buffer = [];
      const isHttps = url.startsWith('https');
      const handler = isHttps ? https : http;
      log.info('isHttps', isHttps);

      const request = handler.get(url, (response) => {
        if (response.statusCode === 200) {
          response.on('data', (chunk) => {
            log.info('[download] request chunk', chunk);
            buffer.push(chunk);
          });
          response.on('end', () => {
            log.info('[download] request end');
            const result = Buffer.concat(buffer);
            log.info('[download] result length', result.length);
            resolve(result);
          });
        } else {
          reject(
            `Server responded with ${response.statusCode}: ${response.statusMessage}`
          );
        }
      });

      request.on('error', (err) => {
        reject(err.message);
      });
    } catch (e) {
      log.error(e);
      reject('Promise error');
    }
  });
}

export function downloadToFile(url: string, dest: string) {
  log.info('[Download]', url, dest);
  const fileName = 'torrents/' + dest + '.torrent';

  return new Promise((resolve, reject) => {
    const errHandler = (err: string) => {
      unlink(fileName, () => {}); // Delete temp file
      reject(err);
    };
    try {
      const file = createWriteStream(fileName, {
        flags: 'wx',
      });
      const isHttps = url.startsWith('https');
      const handler = isHttps ? https : http;
      log.info('isHttps', isHttps);

      const request = handler.get(url, (response) => {
        if (response.statusCode === 200) {
          response.pipe(file);
        } else {
          file.close();
          errHandler(
            `Server responded with ${response.statusCode}: ${response.statusMessage}`
          );
        }
      });

      request.on('error', (err) => {
        file.close();
        errHandler(err.message);
      });

      file.on('finish', () => {
        file.close();
        log.info('[Download] got file ', dest, file);
        resolve(fileName);
      });

      file.on('error', (err) => {
        file.close();
        if (err.code === 'EEXIST') {
          reject('File already exists');
        } else {
          errHandler(err.message);
        }
      });
    } catch (e) {
      log.error(e);
      reject('Promise error');
    }
  });
}

export const readFileAsync = (path: string): Promise<string> => {
  return new Promise(function (resolve, reject) {
    readFile(path, (err, data) => {
      if (!data) reject(err);
      data ? resolve(data.toString()) : resolve(null);
    });
  });
};

export const writeFileAsync = (
  path: string,
  data: string | Buffer
): Promise<void> => {
  try {
    const buffer =
      typeof data === 'string' ? new Uint8Array(Buffer.from(data)) : data;
    return writeFile(path, buffer);
  } catch (err) {
    log.error(err);
  }
  // return new Promise(function (resolve, reject) {
  //   writeFile(path, (err, data) => {
  //     if (!data) reject(err);
  //     data ? resolve(data.toString()) : resolve(null);
  //   });
  // });
};

export function createSuccessResponse(data: any): httpResponse {
  return {
    status: true,
    data: data,
  };
}

export function createErrorResponse(data: any): httpResponse {
  return {
    status: false,
    data: data,
  };
}

const convert = (from?, to?) => (str) => Buffer.from(str, from).toString(to);
export const utf8ToHex = convert('utf8', 'hex');
export const hexToUtf8 = convert('hex', 'utf8');
export const btoa = convert(undefined, 'base64');
export const atob = convert('base64');
