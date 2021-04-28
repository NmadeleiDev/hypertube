import { makeStyles } from '@material-ui/core';
import React, { useRef } from 'react';

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

  const url = id ? `/api/loader/${id}` : src;

  return (
    <div>
      <video
        // poster="poster.png"
        controls
        ref={videoRef}
        className={styles.root}
      >
        <source src={url}></source>
        {/* <source src="video.mp4"></source>
        <source src="video.webm"></source>
        <source src="video.ogv"></source> */}
        <p>This is fallback content</p>
      </video>
    </div>
  );
};

export default NativePlayer;
