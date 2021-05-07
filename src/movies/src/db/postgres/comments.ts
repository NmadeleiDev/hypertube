import log from '../../logger/logger';
import { IComment, IFrontComment } from '../../model/model';
import { query } from './postgres';
const POSTGRES_SCHEME = process.env.POSTGRES_SCHEME || 'hypertube';
const POSTGRES_USERS_SCHEME =
  process.env.POSPOSTGRES_USERS_SCHEMETGRES_SCHEME || 'public';

export const selectCommentById = async (id: number): Promise<IFrontComment> => {
  log.debug('[selectCommentById]', id);
  if (!id) throw new Error('comment id is missing');
  try {
    const res = await query(
      `SELECT
      id,
      movieid,
      text,
      u.image_body AS avatar,
      time,
      u.username AS username
    FROM
    ${POSTGRES_SCHEME}.comments c
    JOIN ${POSTGRES_USERS_SCHEME}.users u ON
      u.user_id = c.userid
      OR u.user_42_id = c.userid
    WHERE
      c.id = $1;`,
      [id]
    );
    log.info(res);
    return res.rows[0] as IFrontComment;
  } catch (e) {
    log.error(e);
  }
};
export const selectCommentsByMovieID = async (
  id: string,
  limit: number = 20,
  offset: number = 0
): Promise<IFrontComment[]> => {
  log.debug('[selectCommentsByMovieID] movieid: ', id);
  if (!id) throw new Error('Movie id is missing');
  try {
    const res = await query(
      `SELECT
        id as commentid,
        movieid,
        text,
        u.image_body AS avatar,
        time,
        u.username AS username
      FROM
      ${POSTGRES_SCHEME}.comments c
      JOIN ${POSTGRES_USERS_SCHEME}.users u ON
        u.user_id = c.userid
        OR u.user_42_id = c.userid
      WHERE
        c.movieid = $1
      LIMIT $2 OFFSET $3;`,
      [id, limit, offset]
    );
    log.info(res);
    return res.rows as IFrontComment[];
  } catch (e) {
    log.error(e);
  }
};

export const insertComment = (comment: IComment) => {
  log.debug('[insertComment]', comment);
  const { userid, movieid, text } = comment;
  const time = comment.time ? new Date(comment.time) : new Date();
  if (!userid || !movieid || !text) throw new Error('userid is missing');
  else if (!comment.movieid) throw new Error('movieid is missing');
  try {
    return query(
      `INSERT INTO ${POSTGRES_SCHEME}.comments(userid, movieid, text, time)
      VALUES ($1, $2, $3, $4) RETURNING *`,
      [userid, movieid, text, time]
    );
  } catch (e) {
    log.debug(comment);
    log.error(e);
    return null;
  }
};

export const updateComment = (comment: IComment) => {
  log.debug('[updateComment]', comment);
  const { id, text } = comment;
  const time = comment.time ? new Date(comment.time) : new Date();
  try {
    if (!id) throw new Error('comment ID id missing');
    return query(
      `UPDATE ${POSTGRES_SCHEME}.comments SET(text, time) = ($2, $3) WHERE id=$1 RETURNING *`,
      [id, text, time]
    );
  } catch (e) {
    log.debug(comment);
    log.error(e);
    return null;
  }
};

export const deleteComment = (id: number) => {
  log.debug('[deleteComment]', id);
  try {
    if (!id) throw new Error('comment ID id missing');
    return query(
      `DELETE FROM ${POSTGRES_SCHEME}.comments WHERE id=$1 RETURNING id`,
      [id]
    );
  } catch (e) {
    log.error(e);
    return null;
  }
};
