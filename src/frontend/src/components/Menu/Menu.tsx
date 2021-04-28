import React from 'react';
import {
  createStyles,
  List,
  ListItem,
  makeStyles,
  Theme,
} from '@material-ui/core';
import { NavLink, useLocation } from 'react-router-dom';
import { secondaryColor, theme } from '../../theme';
import cn from 'classnames';

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'flex-end',
      marginRight: '1rem',
    },
    Item: {
      width: 'unset',
      padding: '0.5rem 1rem',
      fontSize: '1.5rem',
      border: '1px solid transparent',
      transition: '0.3s',
      borderRadius: 4,
      '&:hover': {
        border: '1px solid #ccc',
      },
    },
    Link: {
      color: 'inherit',
      textDecoration: 'none',
    },
    Active: {
      color: secondaryColor.main,
    },
  })
);

const Menu: React.FC = () => {
  const location = useLocation();
  const classes = useStyles(theme);
  const items = [
    {
      id: 0,
      text: 'Popular',
      path: 'popular',
    },
    {
      id: 1,
      text: 'Series',
      path: 'series',
    },
    {
      id: 2,
      text: 'Genres',
      path: 'genres',
    },
  ];

  const getClassName = (path: string): string => {
    return location.pathname.search(path) !== -1
      ? cn(classes.Link, classes.Active)
      : classes.Link;
  };
  return (
    <List className={classes.root}>
      {items.map(({ id, text, path }) => (
        <ListItem key={id} className={classes.Item}>
          <NavLink to={path} className={getClassName(path)}>
            {text}
          </NavLink>
        </ListItem>
      ))}
    </List>
  );
};

export default Menu;
