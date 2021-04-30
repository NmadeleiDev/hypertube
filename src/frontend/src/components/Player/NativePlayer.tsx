import { makeStyles } from '@material-ui/core';
import React, { useRef } from 'react';
import { useSelector } from 'react-redux';
import { RootState } from '../../store/rootReducer';

interface IdProps {
  id: string;
  src?: string;
}

interface SrcProps {
  id?: string;
  src: string;
}

const useSyles = makeStyles({
  root: {
    width: '100%',
    height: '100%',
  },
});

type Props = IdProps | SrcProps;

const NativePlayer = ({ src, id }: Props) => {
  const styles = useSyles();
  const videoRef = useRef<HTMLVideoElement>(null);
  const movie = useSelector((state: RootState) =>
    state.movies.movies.find((movie) => movie.en.id === id)
  );
  const poster = movie?.en.img;

  const url = id ? `/api/storage/load/${id}` : src;

  return (
    <div>
      <video poster={poster} controls ref={videoRef} className={styles.root}>
        <source src={url}></source>
        <p>Cannot play video</p>
      </video>
    </div>
  );
};

export default NativePlayer;
