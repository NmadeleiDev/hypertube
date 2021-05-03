import { makeStyles } from '@material-ui/core';
import React, { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import ReactPlayer from 'react-player';
import screenfull from 'screenfull';
import PlayerControls from './PlayerControls';

const useStyles = makeStyles({
  playerWrapper: {
    position: 'relative',
    width: '100%',
  },
  Error: {
    width: '100%',
    height: '100%',
    padding: 10,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    background: '#ddd',
    borderRadius: 5,
  },
});

interface Props {
  title: string;
  id?: string | number;
  src?: string;
}

const format = (seconds: number) => {
  if (isNaN(seconds)) {
    return '00:00';
  }
  const date = new Date(seconds * 1000);
  const hh = date.getUTCHours();
  const mm = date.getUTCMinutes();
  const ss = date.getUTCSeconds().toString().padStart(2, '0');
  if (hh) {
    return `${hh}:${mm.toString().padStart(2, '0')}:${ss}`;
  }
  return `${mm}:${ss}`;
};

function Player({ title, id, src }: Props) {
  const classes = useStyles();
  const [state, setState] = useState({
    playing: false,
    muted: false,
    volume: 1,
    playbackRate: 1.0,
    played: 0,
    seeking: false,
    controlsVisible: true,
    error: false,
  });
  const [timeDisplayFormat, setTimeDisplayFormat] = useState('normal');
  const {
    playing,
    muted,
    volume,
    playbackRate,
    played,
    seeking,
    controlsVisible,
    error,
  } = state;
  const playerRef = useRef<ReactPlayer>(null);
  const playerContainerRef = useRef<HTMLDivElement>(null);
  const controlsRef = useRef<HTMLDivElement>(null);
  const timeoutId = useRef<NodeJS.Timeout>();
  const videoUrl = id ? `/api/storage/load/${id}/video` : src;
  const { t } = useTranslation();

  useEffect(() => {
    if (!playerRef || !playerRef.current) return;
    console.log('[Player] useEffect');
    const player = playerRef.current.getInternalPlayer();
    console.log(player);
    if (!player) return;
    const playListener = player.addEventListener('play', () => {
      console.log('play');
    });
    const emptiedListener = player.addEventListener('emptied', () => {
      console.log('emptied');
    });
    const durationChangeListener = player.addEventListener(
      'durationchange',
      () => {
        console.log('durationchange');
      }
    );
    const endedListener = player.addEventListener('ended', () => {
      console.log('ended');
    });
    const stalledListener = player.addEventListener('stalled', () => {
      console.log('stalled');
    });
    const suspendListener = player.addEventListener('suspend', () => {
      console.log('suspend');
    });
    const waitingListener = player.addEventListener('waiting', () => {
      console.log('waiting');
    });
    return () => {
      player.removeEventListener('play', playListener);
      player.removeEventListener('emptied', emptiedListener);
      player.removeEventListener('ended', endedListener);
      player.removeEventListener('durationchange', durationChangeListener);
      player.removeEventListener('stalled', stalledListener);
      player.removeEventListener('suspend', suspendListener);
      player.removeEventListener('waiting', waitingListener);
    };
  }, [playerRef]);

  const handlePlayPause = () => {
    setState({ ...state, playing: !state.playing });
  };
  const handleRewind = () => {
    if (playerRef === null || playerRef.current === null) return;
    playerRef.current.seekTo(playerRef.current.getCurrentTime() - 10);
  };
  const handleFastForward = () => {
    if (playerRef === null || playerRef.current === null) return;
    playerRef.current.seekTo(playerRef.current.getCurrentTime() + 10);
  };
  const handleMute = () => {
    setState({ ...state, muted: !state.muted });
  };
  const handleVolumeChange = (
    _e: React.ChangeEvent<{}>,
    newValue: number | number[]
  ) => {
    const value = typeof newValue === 'number' ? newValue : newValue[0];
    setState({
      ...state,
      volume: value / 100,
      muted: value === 0,
    });
  };
  const handleVolumeSeekUp = (
    _e: React.ChangeEvent<{}>,
    newValue: number | number[]
  ) => {
    const value = typeof newValue === 'number' ? newValue : newValue[0];
    setState({
      ...state,
      volume: value / 100,
      muted: value === 0,
    });
  };
  const handlePlaybackRateChange = (rate: number) => {
    setState({ ...state, playbackRate: rate });
  };
  const toggleFullScreen = () => {
    if (!playerContainerRef || !playerContainerRef.current) return;
    if (screenfull && screenfull.isEnabled && playerContainerRef !== null) {
      screenfull.toggle(playerContainerRef.current);
    }
  };

  const handleProgress = ({ played }: { played: number }) => {
    console.log('handleProgress', played);
    if (!controlsRef || !controlsRef.current) return;
    timeoutId.current = setTimeout(() => {
      if (!controlsRef || !controlsRef.current) return;
      controlsRef.current.style.visibility = 'hidden';
    }, 1000);
    if (!state.seeking) {
      setState({ ...state, played });
    }
  };
  const handleSeek = (value: number | number[]) => {
    if (!playerRef || !playerRef.current) return;
    const played = typeof value === 'number' ? value / 100 : value[0] / 100;
    if (played === state.played) return;
    setState({ ...state, played });
    playerRef.current.seekTo(played, 'fraction');
  };
  const handleSeekMouseDown = () => {
    setState({ ...state, seeking: true });
  };
  const handleSeekMouseUp = () => {
    setState({ ...state, seeking: false });
  };
  const handleChangeDisplayFormat = () => {
    setTimeDisplayFormat(
      timeDisplayFormat === 'normal' ? 'remaining' : 'normal'
    );
  };
  const handleMouseMove = () => {
    if (!controlsRef || !controlsRef.current) return;
    controlsRef.current.style.visibility = 'visible';
    if (timeoutId.current) clearTimeout(timeoutId.current);
  };
  const handleError = (e: Error) => {
    console.log('handleError', e);

    setState((prev) => ({ ...prev, played: prev.played + 0.001 }));
  };

  const currentTime =
    !playerRef || !playerRef.current ? 0 : playerRef.current.getCurrentTime();
  const duration =
    !playerRef || !playerRef.current ? 0 : playerRef.current.getDuration();
  const elapsedTime =
    timeDisplayFormat === 'normal'
      ? format(currentTime)
      : `-${format(duration - currentTime)}`;
  const totalDuration = format(duration);

  return error ? (
    <div className={classes.Error}>{t`Cannot load movie`}</div>
  ) : (
    <div
      ref={playerContainerRef}
      className={classes.playerWrapper}
      onMouseMove={handleMouseMove}
    >
      <ReactPlayer
        ref={playerRef}
        width="100%"
        height="100%"
        url={videoUrl}
        muted={muted}
        playing={playing}
        volume={volume}
        controls={false}
        playbackRate={playbackRate}
        onProgress={handleProgress}
        onError={handleError}
        onBuffer={() => console.log('onBuffer')}
        onBufferEnd={() => console.log('onBufferEnd')}
        // config={{
        //   file: {
        //     tracks: [
        //       {
        //         kind: 'captions',
        //         src: `/api/storage/load/${id}/srt`,
        //         srcLang: 'en',
        //         label: 'en',
        //         default: true,
        //       },
        //     ],
        //   },
        // }}
      />
      {controlsVisible && (
        <PlayerControls
          ref={controlsRef}
          onPlayPause={handlePlayPause}
          onMuted={handleMute}
          playing={playing}
          muted={muted}
          volume={volume}
          played={played}
          seeking={seeking}
          onRewind={handleRewind}
          onFastForward={handleFastForward}
          onVolumeChange={handleVolumeChange}
          onVolumeSeekUp={handleVolumeSeekUp}
          playbackRate={playbackRate}
          onPlaybackRateChange={handlePlaybackRateChange}
          onToggleFullScreen={toggleFullScreen}
          onSeek={handleSeek}
          onSeekMouseDown={handleSeekMouseDown}
          onSeekMouseUp={handleSeekMouseUp}
          elapsedTime={elapsedTime}
          totalDuration={totalDuration}
          onChangeDisplayFormat={handleChangeDisplayFormat}
          title={title}
        />
      )}
    </div>
  );
}

export default Player;
