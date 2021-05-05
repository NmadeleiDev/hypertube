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
  Controls: {},
  Progress: {},
});

type Props = IdProps | SrcProps;

const NativePlayer = ({ src, id }: Props) => {
  const styles = useSyles();
  const videoRef = useRef<HTMLVideoElement>(null);
  const controlsRef = useRef<HTMLDivElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
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

  useEffect(() => {
    if (!videoRef || !videoRef.current) return;
    if (!controlsRef || !controlsRef.current) return;
    videoRef.current.controls = false;
    controlsRef.current.style.display = 'block';
  });

  return (
    <div ref={containerRef}>
      <video
        id="video"
        controls
        ref={videoRef}
        preload="metadata"
        className={styles.root}
      >
        <source src={url} type="video/mp4" />
        <track
          label="English"
          kind="subtitles"
          srcLang="en"
          src={`/api/storage/load/${id}/srt`}
          default
        />
        <p>Cannot play video</p>
      </video>
      <div
        id="video-controls"
        ref={controlsRef}
        className={styles.Controls}
        data-state="hidden"
      >
        <button id="playpause" type="button" data-state="play">
          Play/Pause
        </button>
        <button id="stop" type="button" data-state="stop">
          Stop
        </button>
        <div className={styles.Progress}>
          <progress id="progress" value="0" max="100">
            <span id="progress-bar"></span>
          </progress>
        </div>
        <button id="mute" type="button" data-state="mute">
          Mute/Unmute
        </button>
        <button id="volinc" type="button" data-state="volup">
          Vol+
        </button>
        <button id="voldec" type="button" data-state="voldown">
          Vol-
        </button>
        <button id="fs" type="button" data-state="go-fullscreen">
          Fullscreen
        </button>
        <button id="subtitles" type="button" data-state="subtitles">
          CC
        </button>
      </div>
    </div>
  );
};

export default NativePlayer;
