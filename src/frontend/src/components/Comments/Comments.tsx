import React, { useEffect, useRef } from 'react';
import { Grid, Typography } from '@material-ui/core';
import Comment from './Comment/Comment';
import CategoryHeader from '../CategoryHeader/CategoryHeader';
import SendComment from './SendComment/SendComment';
import { loadComments } from '../../store/features/MoviesSlice';
import { debouncedDetectBottomLine } from '../../utils';
import { useTranslation } from 'react-i18next';
import { useAppDispatch } from '../../store/store';
import { IMovie } from '../../models/MovieInfo';
import { useSelector } from 'react-redux';
import { RootState } from '../../store/rootReducer';

const LIMIT = 5;
export interface CommentsProps {
  movie: IMovie;
}

const Comments: React.FC<CommentsProps> = ({ movie }) => {
  const gridRef = useRef<HTMLDivElement>(null);
  const dispatch = useAppDispatch();
  const { t } = useTranslation();
  const offsetRef = useRef(0);
  const { isAuth } = useSelector((state: RootState) => state.user);

  // load first batch of comments initially
  // if there's something left to load
  useEffect(() => {
    if (!movie?.info.maxComments) return;
    if (
      !movie.info.comments ||
      (movie.info.comments &&
        movie.info.comments.length < movie.info.maxComments)
    ) {
      dispatch(
        loadComments(
          { movieId: movie.id, offset: movie.info.comments?.length },
          (length) => (offsetRef.current += length)
        )
      );
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    if (!movie || !movie.info.maxComments) return;
    const endOfCommentsCallback = () => {
      if (
        movie.info.comments?.length &&
        movie.info.comments?.length < movie.info.maxComments
      ) {
        dispatch(
          loadComments({
            movieId: movie.id,
            limit: LIMIT,
            offset: movie.info.comments?.length || 0,
          })
        );
      }
    };
    const trackScrolling = debouncedDetectBottomLine(
      gridRef.current,
      endOfCommentsCallback
    );
    window.addEventListener('scroll', trackScrolling);
    return () => window.removeEventListener('scroll', trackScrolling);
  }, [
    movie,
    movie.info.comments?.length,
    movie.info.maxComments,
    movie.id,
    dispatch,
  ]);

  if (!movie) return null;
  console.log('[Comments]', movie.info.comments);
  const content = movie.info.comments?.length ? (
    movie.info.comments.map((comment) => (
      <Comment {...comment} key={comment.id} />
    ))
  ) : (
    <Typography variant="body1" style={{ fontSize: '1.3rem' }}>
      {t`No comments`}
    </Typography>
  );

  return (
    <Grid ref={gridRef} style={{ marginTop: 10 }} container direction="column">
      <CategoryHeader text={t`Comments`} />
      {isAuth && <SendComment />}
      {content}
    </Grid>
  );
};

export default Comments;
