import { Button, ButtonGroup, Grid, makeStyles } from '@material-ui/core';
import React from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import { theme } from '../../theme';
import { useTranslation } from 'react-i18next';

export interface ILink {
  id: number;
  to: string;
  name: string;
  highlight?: boolean;
}

export interface NavProps {
  links?: ILink[];
}

const defaultLinks: ILink[] = [
  { id: 0, to: '/', name: 'navHome' },
  { id: 1, to: '/films', name: 'navFilms' },
  { id: 2, to: '/series', name: 'navSeries' },
  { id: 3, to: '/new', name: 'navNew' },
];

const useStyles = makeStyles({
  root: {
    [theme.breakpoints.down('xs')]: {
      display: 'none',
    },
  },
  Link: {
    textDecoration: 'none',
    color: 'inherit',
  },
});

const Nav: React.FC<NavProps> = ({ links = defaultLinks }) => {
  const classes = useStyles();
  const location = useLocation();
  const { t } = useTranslation();

  return (
    <Grid className={classes.root} container justify="center">
      <ButtonGroup style={{ marginBottom: 20 }}>
        {links.map(({ id, to, highlight, name }) => (
          <Button
            key={id}
            component={NavLink}
            to={to}
            variant={location.pathname === to ? 'contained' : 'outlined'}
            color={highlight ? 'secondary' : 'primary'}
          >
            {t(name)}
          </Button>
        ))}
      </ButtonGroup>
    </Grid>
  );
};

export default Nav;
