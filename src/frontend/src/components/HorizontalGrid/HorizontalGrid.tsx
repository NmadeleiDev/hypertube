import { Grid, makeStyles } from '@material-ui/core';
import React from 'react';
import { useTranslation } from 'react-i18next';
import CategoryHeader from '../CategoryHeader/CategoryHeader';

interface Props {
  name: string;
  sources?: string[];
  type: 'photo' | 'video';
}
interface GridItemProps {
  name: string;
  src: string;
  type: 'photo' | 'video';
  key?: string | number;
}

const useStyles = makeStyles((theme) => ({
  root: {
    padding: '10px 0',
  },
  gridList: {
    flexWrap: 'nowrap',
    transform: 'translateZ(0)',
  },
  HorizontalScroll: {
    overflowX: 'auto',
  },
  Image: {
    height: '10rem',
    marginRight: 10,
  },
}));

const GridItem: React.FC<GridItemProps> = ({
  src,
  name,
  type,
}: GridItemProps) => {
  const classes = useStyles();

  switch (type) {
    case 'photo':
      return <img src={src} alt={name} className={classes.Image} />;
    case 'video':
      return <video controls src={src} className={classes.Image} />;
    default:
      return null;
  }
};

const HorizontalGrid = ({ sources, name, type }: Props) => {
  const classes = useStyles();
  const { t } = useTranslation();

  return (
    <Grid container direction="column" className={classes.root}>
      <CategoryHeader text={t(`${type}s`)} />
      <Grid container wrap="nowrap" className={classes.HorizontalScroll}>
        {sources?.length
          ? sources.map((item: string, index: number) => (
              <GridItem key={item + index} src={item} name={name} type={type} />
            ))
          : t`No info`}
      </Grid>
    </Grid>
  );
};

export default HorizontalGrid;
