import { Container, createStyles, makeStyles, Theme } from '@material-ui/core';
import React, { useState } from 'react';
import { throttle } from 'lodash';
import { NavLink, useLocation } from 'react-router-dom';
import cn from 'classnames';
import { secondaryColor, theme } from '../../theme';
import { useAppDispatch } from '../../store/store';
import {
  loadMovies,
  resetEndOfMovies,
  setMovies,
} from '../../store/features/MoviesSlice';
import { LIMIT } from '../..';
import { resetFilterState } from '../../store/features/FilterSlice';

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
      padding: 0,
      position: 'relative',
      borderRadius: 5,
      overflow: 'hidden',
    },
    Scroller: {
      display: 'inline-flex',
      flexWrap: 'nowrap',
      overflowX: 'scroll',
      '&::-webkit-scrollbar': {
        display: 'none',
      },
    },
    ShadowLeft: {
      position: 'absolute',
      height: '100%',
      left: 0,
      width: '50px',
      background: 'linear-gradient(90deg, #00000033, transparent)',
    },
    ShadowRight: {
      position: 'absolute',
      height: '100%',
      right: 0,
      width: '50px',
      background: 'linear-gradient(-90deg, #00000033, transparent)',
    },
    Link: {
      color: theme.palette.grey[700],
      fontSize: '1.5rem',
      textDecoration: 'none',
      flex: '1 0 50px',
      margin: 10,
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      border: '3px solid transparent',
      borderRadius: '50%',
      '&:before': {
        content: "''",
        paddingTop: '100%',
        float: 'left',
      },
    },
    Active: {
      borderColor: secondaryColor.main,
    },
  })
);

const AlphabetNav = () => {
  const location = useLocation();
  const classes = useStyles(theme);
  const [shadow, setShadow] = useState({
    shadowLeft: false,
    shadowRight: true,
  });
  const letters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
  const dispatch = useAppDispatch();

  const handleLinkClick = async (
    e: React.SyntheticEvent<HTMLAnchorElement>
  ) => {
    const letter = e.currentTarget.innerText;
    const isByNamePage = window.location.href.match(/byname/);
    // check if we were on one of 'byname' pages
    if (isByNamePage) {
      // check if user clicks the same letter, we're in
      const currentLetter = window.location.href.split('/').pop();
      // no need to load same movies twice
      if (currentLetter && letter === currentLetter) return;
    }
    dispatch(resetEndOfMovies());
    dispatch(setMovies({ byName: [] }));
    dispatch(resetFilterState());
    dispatch(loadMovies({ filter: { letter, limit: LIMIT } }));
  };

  const scrollHandler = throttle((e: React.SyntheticEvent<HTMLDivElement>) => {
    e.stopPropagation();
    const target = e.target as HTMLDivElement;
    if (!target) return;
    if (target.scrollLeft === 0) {
      setShadow({ ...shadow, shadowLeft: false });
    } else if (target.scrollLeft > 0) {
      !shadow.shadowLeft && setShadow({ ...shadow, shadowLeft: true });
    }
    if (target.scrollWidth - target.scrollLeft < target.offsetWidth + 1) {
      setShadow({ ...shadow, shadowRight: false });
    } else {
      !shadow.shadowRight && setShadow({ ...shadow, shadowRight: true });
    }
  }, 100);

  const horizontalScrollHandler = (e: React.WheelEvent<HTMLDivElement>) => {
    e.stopPropagation();
    e.currentTarget.scrollTo({
      top: 0,
      left: e.currentTarget.scrollLeft + e.deltaY * 4,
      behavior: 'smooth',
    });
  };

  const getItemClassName = (path: string): string => {
    return location.pathname.search(path) !== -1
      ? cn(classes.Link, classes.Active)
      : classes.Link;
  };

  return (
    <Container className={classes.root}>
      {shadow.shadowLeft && <div className={classes.ShadowLeft} />}
      {shadow.shadowRight && <div className={classes.ShadowRight} />}
      <Container
        onScroll={scrollHandler}
        onWheel={horizontalScrollHandler}
        className={classes.Scroller}
      >
        {letters.split('').map((letter) => (
          <NavLink
            to={`/byname/${letter}`}
            key={letter}
            onClick={handleLinkClick}
            className={getItemClassName(`/byname/${letter}`)}
          >
            {letter}
          </NavLink>
        ))}
        <div className={classes.Link} />
      </Container>
    </Container>
  );
};

export default AlphabetNav;
