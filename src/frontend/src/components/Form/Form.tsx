import React from 'react';
import {
  Button,
  ButtonProps,
  List,
  ListItem,
  makeStyles,
} from '@material-ui/core';
import Input, { IInputProps } from '../Input/Input';
import cn from 'classnames';
import { NavLink } from 'react-router-dom';

const useStyles = makeStyles({
  root: {
    width: '100%',
  },
  Button: {
    width: '100%',
  },
  List: {
    width: '100%',
  },
  Buttons: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
  },
  Link: {
    textDecoration: 'none',
    color: 'inherit',
    width: '100%',
  },
  margin: {
    marginBottom: 10,
  },
  padding: {
    padding: '5px 10px',
  },
});

export interface IButtonProps extends ButtonProps {
  text: string;
  type?: 'submit';
  to?: string;
  disabled?: boolean;
  variant?: 'contained' | 'outlined' | 'text';
  onClick: (e: React.MouseEvent<HTMLButtonElement>) => {};
}

export interface IFormProps {
  inputs: IInputProps[];
  buttons: IButtonProps[];
  className?: string;
}

const Form: React.FC<IFormProps> = ({
  inputs,
  buttons,
  className,
  children,
}) => {
  const classes = useStyles();

  return (
    <form className={cn(className, classes.root)}>
      <List className={classes.List}>
        {inputs.map((el) => (
          <ListItem key={el.name || el.type}>
            <Input {...el} className={classes.margin} />
          </ListItem>
        ))}
        {buttons && (
          <ListItem className={classes.Buttons}>
            {buttons.map(
              ({ variant, type, text, onClick, to, ...others }) =>
                to ? (
                  <NavLink
                    to={to}
                    key={text}
                    className={cn(classes.Link, classes.margin)}
                  >
                    <Button
                      className={classes.Button}
                      variant={variant}
                      onClick={onClick}
                      {...others}
                    >
                      {text}
                    </Button>
                  </NavLink>
                ) : (
                  <Button
                    key={text}
                    className={cn(classes.Button, classes.margin)}
                    variant={variant}
                    type={type}
                    onClick={onClick}
                    {...others}
                  >
                    {text}
                  </Button>
                )
              // )
            )}
          </ListItem>
        )}
        {children}
      </List>
    </form>
  );
};

export default Form;
