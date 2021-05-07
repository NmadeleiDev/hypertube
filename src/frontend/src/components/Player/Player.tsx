import { makeStyles } from '@material-ui/core';
import axios from 'axios';
import React, { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import ReactPlayer from 'react-player';
import { TrackProps } from 'react-player/file';
import screenfull from 'screenfull';
import { useToast } from '../../hooks/useToast';
import { getSearchParam } from '../../utils';
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
  id: string | number;
  tracksProps: TrackProps[];
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

function Player({ title, id, tracksProps }: Props) {
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
    showSubtitles: true,
    loading: false,
  });
  const [track, setTrack] = useState<number>(0);
  const [player, setPlayer] = React.useState<Record<string, any>>();
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
    showSubtitles,
    loading,
  } = state;
  const playerRef = useRef<ReactPlayer>(null);
  const playerContainerRef = useRef<HTMLDivElement>(null);
  const controlsRef = useRef<HTMLDivElement>(null);
  const timeoutId = useRef<NodeJS.Timeout>();
  const videoUrl = `/api/storage/load/${id}/video`;
  const { toast } = useToast();
  const { t } = useTranslation();

  /**
   * set played to the value, defined in params
   * could be used if there is an error -
   * then reload component and set current played position
   * to the value, set in parameters
   */
  useEffect(() => {
    const params = getSearchParam();
    console.log(params);
    if (params?.time) {
      const played = parseFloat(params.time);
      setState((state) => ({ ...state, played }));
    }
  }, []);

  // save internal player
  useEffect(() => {
    if (!playerRef || !playerRef.current) return;
    console.log('[Player] <save internal player> useEffect');
    const player = playerRef.current.getInternalPlayer();
    setPlayer(player);
    console.log('[Player] <save internal player> player', player);
  }, [playerRef]);

  // add listners
  useEffect(() => {
    if (!player) return;
    console.log('[Player] <add listners> useEffect');
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
  }, [player]);

  const handleSwitchSubtitles = (index: number) => {
    setTrack(index);
  };

  const handleChangeSubtitles = () => {
    if (playerRef === null || playerRef.current === null) return;
    if (!player) {
      const player = playerRef.current.getInternalPlayer();
      setPlayer(() => player);
    }
    console.log('[Player] <handleChangeSubtitles> player', player);
    const tracks = player?.textTracks;
    console.log('[Player] <handleChangeSubtitles> tracks', tracks);
    if (!tracks || !tracks.length) return;
    const showing = tracks[track].mode === 'showing';
    console.log(
      '[Player] <handleChangeSubtitles> track, showing, showSubtitles',
      track,
      showing,
      showSubtitles
    );
    tracks[track].mode = showing ? 'disabled' : 'showing';
    setState({ ...state, showSubtitles: !showing });
  };

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
    setState({ ...state, loading: false });
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
    toast({ text: t`movieError` }, 'error', { autoHideDuration: 7000 });
  };
  const handleBuffer = () => {
    console.log('onBuffer');
    setState({ ...state, loading: true });
  };
  const handleBufferEnd = () => {
    console.log('onBufferEnd');
    setState({ ...state, loading: false });
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
        onBuffer={handleBuffer}
        onBufferEnd={handleBufferEnd}
        onReady={handleBufferEnd}
        config={{ file: { tracks: [...tracksProps] } }}
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
          loading={loading}
          showSubtitles={showSubtitles}
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
          onChangeSubtitles={handleChangeSubtitles}
          onSwitchSubtitles={handleSwitchSubtitles}
          onChangeDisplayFormat={handleChangeDisplayFormat}
          title={title}
          tracks={tracksProps}
        />
      )}
    </div>
  );
}

export default Player;
