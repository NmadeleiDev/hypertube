import { makeStyles } from '@material-ui/core';
import React, { useEffect, useRef } from 'react';
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
    maxHeight: '400px',
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

  useEffect(() => {
    if (!videoRef || !videoRef.current) return;
    const video = videoRef.current;
    const handleProgress = () => {
      console.log('[handleProgress]');
      const buffered = video.buffered;
      const currentTime = video.currentTime;
      console.log('buffered', buffered);
      console.log('currentTime', currentTime);
    };
    video.addEventListener('progress', handleProgress);
    return () => {
      video.removeEventListener('progress', handleProgress);
    };
  });

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
