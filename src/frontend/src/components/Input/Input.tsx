import React, { useEffect, useState } from 'react';
import {
  makeStyles,
  TextField,
  // MenuItem,
  // Box,
  // Slider,
  // Typography,
  // Grid,
  // FormControl,
  // FormLabel,
  // RadioGroup,
  // FormControlLabel,
  // Radio,
} from '@material-ui/core';

const useStyles = makeStyles({
  ButtonListItem: {
    justifyContent: 'center',
  },
  Input: {
    width: '100%',
  },
  SmallSelector: {
    minWidth: '100px',
    padding: '10px',
  },
  FormControl: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'flex-start',
  },
  RadioGroup: {
    width: '100%',
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-evenly',
    alignItems: 'center',
  },
});

export interface IRule {
  minLength: number;
  maxLength: number;
  regex?: RegExp;
  checkFn?: () => boolean;
}

export interface IInputProps {
  name: string;
  type?: 'text' | 'email' | 'password';
  label?: string;
  value: string;
  size?: 'medium' | 'small';
  fullWidth?: boolean;
  required?: boolean;
  values?: string[];
  placeholder?: string;
  ignoreUntouched?: boolean;
  onValidate?: (arg0: boolean) => {};
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => {};
  className: string;
  rules?: {
    helperText: string;
    rule: IRule;
  };
}

// const getDate = (timestamp: number): string => {
//   const now = new Date();
//   const date =
//     timestamp > now.getTime() + 1000 * 60 * 60 * 24 ? now : new Date(timestamp);
//   const year = date.getFullYear();
//   let yearStr = '';
//   if (year < 10) yearStr = '000' + year;
//   else if (year < 100) yearStr = '00' + year;
//   else if (year < 1000) yearStr = '0' + year;
//   const month = date.getMonth() + 1;
//   const monthStr = month < 10 ? '0' + month : '' + month;
//   const day = date.getDate();
//   const dayStr = day < 10 ? '0' + day : '' + day;
//   const formattedDate = `${yearStr}-${monthStr}-${dayStr}`;
//   return formattedDate;
// };

/**
 * checks *value* based on *rule* if input was *touched*
 * @param {string} value checked value
 * @param {IRule} rule rule to match value
 * @param {boolean} touched flag to check if input was already focused
 */
const validate = (value: string, rule: IRule, touched: boolean): boolean => {
  const isRegexValid = rule.regex ? !!value.match(rule.regex) : true;
  const isFnValid = rule.checkFn ? rule.checkFn() : true;
  return touched
    ? value.length >= rule.minLength &&
        value.length <= rule.maxLength &&
        isRegexValid &&
        isFnValid
    : true;
};

const Input = ({
  type = 'text',
  name,
  values,
  value,
  onValidate,
  onChange,
  className,
  rules,
  ignoreUntouched,
  ...rest
}: IInputProps) => {
  const classes = useStyles();
  const [touched, setTouched] = useState(false);
  const [error, setError] = useState(false);
  const [errorText, setErrorText] = useState('');

  const { helperText, rule } = React.useMemo(
    () => (rules ? { ...rules } : { ...{ helperText: '', rule: null } }),
    [rules]
  );

  useEffect(() => {
    let err;
    // check *touched* only if *ignoreUntouched* is false
    err = rule ? !validate(value, rule, ignoreUntouched || touched) : false;
    // update error message
    err ? setErrorText(helperText) : setErrorText('');
    // update inputValid state in Redux
    onValidate &&
      onValidate((ignoreUntouched || touched || !!value.length) && !err);
    if (error !== err) setError(err);
  }, [value, touched, rule, helperText, onValidate, ignoreUntouched, error]);

  switch (type) {
    // case 'number':
    case 'text':
    case 'email':
    case 'password':
      return (
        <TextField
          id={name}
          error={error}
          helperText={errorText}
          type={type}
          value={value}
          name={name}
          className={classes.Input}
          onFocus={() => setTouched(true)}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => onChange(e)}
          {...rest}
        />
      );
    /*
    case 'select':
      return (
        <TextField
          id={name}
          key={name}
          label={label}
          value={value}
          name={name}
          className={classes.Input}
          onChange={(e) => onChange(e)}
          select
          {...rest}
        >
          {values.map((e) => (
            <MenuItem key={e} value={e}>
              {e}
            </MenuItem>
          ))}
        </TextField>
      );
    case 'radio':
      return (
        <FormControl
          className={classes.FormControl}
          key={name}
          component="fieldset"
        >
          <FormLabel component="legend">{capitalize(label)}</FormLabel>
          <RadioGroup
            className={classes.RadioGroup}
            value={value}
            aria-label={name}
            name={name}
            onChange={(e) => onChange(e)}
          >
            {values.map((el) => (
              <FormControlLabel
                key={el}
                value={el}
                control={<Radio checked={value === el} color="primary" />}
                label={capitalize(el)}
              />
            ))}
          </RadioGroup>
        </FormControl>
      );
    case 'date':
      return (
        <TextField
          id={name}
          type="date"
          label="Birth Date"
          name={name}
          className={classes.Input}
          InputLabelProps={{ shrink: true }}
          value={getDate(value)}
          onChange={(e) => onChange(e)}
          {...rest}
        />
      );
    case 'slider':
      return (
        <Box className={classes.Input}>
          <Typography id="range-slider" gutterBottom>
            {label}
          </Typography>
          <Slider
            value={value}
            onChange={(e, value) => onChange(value)}
            valueLabelDisplay="auto"
            aria-labelledby="range-slider"
            {...rest}
          />
        </Box>
      );
    case 'doubleSelector':
      return (
        <Box key={name} className={classes.Input}>
          <Grid container justify="space-evenly" alignItems="center">
            <Grid item xs={4}>
              <TextField
                shrink
                id={'min' + name}
                type="number"
                value={value.minAge}
                label="Minimum age"
                className={classes.SmallSelector}
                onChange={(e) => onChange({ minAge: +e.target.value })}
                {...rest}
              />
            </Grid>
            <Grid item xs={4}>
              <TextField
                shrink
                id={'max' + name}
                type="number"
                value={value.maxAge}
                label="Maximum age"
                className={classes.SmallSelector}
                onChange={(e) => onChange({ maxAge: +e.target.value })}
                {...rest}
              />
            </Grid>
          </Grid>
        </Box>
      );
    case 'autocomplete':
      return (
        <TagInput
          name={name}
          className={classes.Input}
          tags={values}
          onChange={onChange}
          {...rest}
        />
      );

      */
    default:
      return null;
  }
};

export default Input;
