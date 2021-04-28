import {
  Button,
  ButtonProps,
  Grid,
  makeStyles,
  Paper,
  Popover,
} from '@material-ui/core';
import React, { createContext, useRef, useState } from 'react';

export interface IItem {
  text: string;
  name: string;
  onClick?: (e: React.MouseEvent<HTMLButtonElement>) => void;
}
interface Props {
  heroText?: string;
  icon?: JSX.Element;
  items?: IItem[];
  buttonProps?: ButtonProps;
}

const useStyles = makeStyles({
  root: {
    position: 'relative',
  },
  Paper: {
    padding: 20,
  },
  Open: {
    transform: 'rotate(180deg)',
  },
  Item: {
    fontSize: '1.0rem',
    margin: 5,
  },
});

const contextValue = {
  onClose: () => {},
};

export const handlersContext = createContext(contextValue);

const Dropdown: React.FC<Props> = ({
  heroText = null,
  icon,
  items,
  buttonProps,
  children,
}) => {
  const [open, setOpen] = useState(false);
  const anchorEl = useRef<HTMLButtonElement>(null);
  const classes = useStyles();

  const handleButtonClick = (
    e: React.MouseEvent<HTMLButtonElement>,
    onClick: undefined | ((e: React.MouseEvent<HTMLButtonElement>) => void)
  ) => {
    onClick && onClick(e);
    setOpen(false);
  };

  const handleClose = () => setOpen(false);
  const contextValue = {
    onClose: handleClose,
  };

  return (
    <div className={classes.root}>
      <Button
        style={{ margin: '0 10px' }}
        ref={anchorEl}
        variant="text"
        onClick={() => setOpen((open) => !open)}
        endIcon={heroText && icon}
        {...buttonProps}
      >
        {heroText || icon}
      </Button>
      <Popover
        open={open}
        anchorEl={anchorEl.current}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
        transformOrigin={{ vertical: 'top', horizontal: 'center' }}
        onClose={handleClose}
      >
        <Paper className={classes.Paper}>
          <Grid container direction="column" alignItems="center">
            {items &&
              items.map(({ text, name, onClick }: IItem) => (
                <Button
                  key={text}
                  variant="text"
                  name={name}
                  className={classes.Item}
                  onClick={(e) => handleButtonClick(e, onClick)}
                >
                  {text}
                </Button>
              ))}
            <handlersContext.Provider value={contextValue}>
              {children}
            </handlersContext.Provider>
          </Grid>
        </Paper>
      </Popover>
    </div>
  );
};

export default Dropdown;
