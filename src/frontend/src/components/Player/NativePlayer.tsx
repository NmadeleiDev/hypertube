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
    const player = videoRef.current;

    console.log(player);
    if (!player) return;
    const playListener = () => {
      console.log('play');
    };
    const emptiedListener = () => {
      console.log('emptied');
    };
    const durationChangeListener = () => {
      console.log('durationchange');
    };
    const endedListener = () => {
      console.log('ended');
    };
    const stalledListener = () => {
      console.log('stalled');
    };
    const suspendListener = () => {
      console.log('suspend');
    };
    const waitingListener = () => {
      console.log('waiting');
    };
    const handleProgress = () => {
      console.log('[handleProgress]');
      const buffered = player.buffered;
      const currentTime = player.currentTime;
      console.log('buffered', buffered);
      console.log('currentTime', currentTime);
    };
    player.addEventListener('play', playListener);
    player.addEventListener('emptied', emptiedListener);
    player.addEventListener('ended', endedListener);
    player.addEventListener('durationchange', durationChangeListener);
    player.addEventListener('stalled', stalledListener);
    player.addEventListener('suspend', suspendListener);
    player.addEventListener('waiting', waitingListener);
    player.addEventListener('progress', handleProgress);
    return () => {
      player.removeEventListener('play', playListener);
      player.removeEventListener('emptied', emptiedListener);
      player.removeEventListener('ended', endedListener);
      player.removeEventListener('durationchange', durationChangeListener);
      player.removeEventListener('stalled', stalledListener);
      player.removeEventListener('suspend', suspendListener);
      player.removeEventListener('waiting', waitingListener);
      player.removeEventListener('progress', handleProgress);
    };
  }, []);

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
