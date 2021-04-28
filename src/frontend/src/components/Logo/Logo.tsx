import { makeStyles, SvgIcon } from '@material-ui/core';
import React from 'react';

interface Props {
  width?: number;
  height?: number;
}

const useStyles = makeStyles({
  root: {
    color: '#035263',
    '&:hover>': {
      '& $_0': {
        animation: `$forward 1s infinite ease-in-out`,
      },
      '& $_1': {
        animation: `$stretch1 1s infinite ease-out`,
      },
      '& $_2': {
        animation: `$stretch2 0.9s infinite ease-out`,
      },
      '& $_3': {
        animation: `$stretch1 1.1s infinite ease-out`,
      },
      '& $_4': {
        animation: `$stretch1 0.8s infinite ease-out`,
      },
    },
  },
  _0: {},
  _1: {},
  _2: {},
  _3: {},
  _4: {},
  '@keyframes stretch1': {
    '0%': {
      opacity: 0.8,
      transform: 'translateX(0) scaleX(1)',
    },
    '10%': {
      opacity: 1,
      transform: 'translateX(-10px) scaleX(1)',
    },
    '70%': {
      opacity: 0.8,
    },
    '80%': {
      opacity: 0.2,
      transform: 'translateX(-80px) scaleX(1.5)',
    },
    '100%': {
      opacity: 0,
    },
  },
  '@keyframes forward': {
    '0%': {
      transform: 'translateX(0)',
    },
    '10%': {
      transform: 'translateX(10px) translateY(3px)',
    },
    '50%': {
      transform: 'translateX(0px) translateY(-3px)',
    },
    '80%': {
      transform: 'translateX(-7px) translateY(4px)',
    },
    '100%': {
      transform: 'translateX(0px) translateY(0)',
    },
  },
  '@keyframes stretch2': {
    '0%': {
      opacity: 0.2,
      transform: 'translateX(0) scaleX(1)',
    },
    '10%': {
      opacity: 1,
      transform: 'translateX(-15px) scaleX(1.2)',
    },
    '80%': {
      opacity: 0.6,
    },
    '90%': {
      opacity: 0.2,
      transform: 'translateX(-40px) scaleX(1.4)',
    },
    '100%': {
      opacity: 0,
    },
  },
});

const Logo = ({ width = 70, height = 50 }: Props) => {
  const classes = useStyles();
  return (
    <SvgIcon
      style={{ width, height }}
      className={classes.root}
      viewBox={'-100 0 584 384'}
    >
      <path
        className={classes._0}
        d="M376.304,178.327l-224-136c-4.944-3-11.12-3.104-16.152-0.272C131.12,44.887,128,50.215,128,55.999v136    c0,8.832,7.168,16,16,16c8.832,0,16-7.168,16-16V84.431l177.168,107.568L160,299.567v-43.568c0-8.832-7.168-16-16-16    c-8.832,0-16,7.168-16,16v72c0,5.784,3.12,11.112,8.152,13.944c2.44,1.376,5.144,2.056,7.848,2.056    c2.88,0,5.752-0.776,8.304-2.328l224-136c4.776-2.896,7.696-8.08,7.696-13.672C384,186.407,381.08,181.223,376.304,178.327z"
      />
      <path
        className={classes._1}
        d="M88,279.999H16c-8.832,0-16,7.168-16,16c0,8.832,7.168,16,16,16h72c8.832,0,16-7.168,16-16    C104,287.167,96.832,279.999,88,279.999z"
      />
      <path
        className={classes._2}
        d="M56,71.999H16c-8.832,0-16,7.168-16,16c0,8.832,7.168,16,16,16h40c8.832,0,16-7.168,16-16    C72,79.167,64.832,71.999,56,71.999z"
      />
      <path
        className={classes._3}
        d="M88,135.999h-8c-8.832,0-16,7.168-16,16c0,8.832,7.168,16,16,16h8c8.832,0,16-7.168,16-16    C104,143.167,96.832,135.999,88,135.999z"
      />
      <path
        className={classes._4}
        d="M72,207.999H48c-8.832,0-16,7.168-16,16c0,8.832,7.168,16,16,16h24c8.832,0,16-7.168,16-16    C88,215.167,80.832,207.999,72,207.999z"
      />
    </SvgIcon>
  );
};

export default Logo;
