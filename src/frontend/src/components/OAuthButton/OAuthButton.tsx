import { Button, CircularProgress, makeStyles } from '@material-ui/core';
import React, { useMemo, useState } from 'react';

interface Props {
  variant: '42' | 'facebook' | 'vk' | 'google';
}

const useStyles = makeStyles({
  OAuthButton: {
    width: '6rem',
  },
});

const oauth = {
  '42': {
    url: 'https://api.intra.42.fr/oauth/authorize',
    queryParams: [
      {
        name: 'client_id',
        value:
          '96975efecfd0e5efee67c9ac4cc350ac9372ae559b2fb8a08feba6841a33fb53',
      },
      {
        name: 'redirect_uri',
        value: 'http://localhost:4000/api/auth/oauth42',
      },
      {
        name: 'scope',
        value: 'public',
      },
      {
        name: 'state',
        value:
          'bdcbe28874ab05962b50430b1466a8ebcbda45ba8c3c1beee600699478ad2a4d',
      },
      {
        name: 'response_type',
        value: 'code',
      },
    ],
  },
  facebook: {
    url: 'https://www.facebook.com/v10.0/dialog/oauth',
    queryParams: [
      {
        name: 'client_id',
        value: '737140593612593',
      },
      {
        name: 'redirect_uri',
        value: 'http://localhost:4000/api/auth/oauthFb',
      },
      {
        name: 'scope',
        value: 'email+public_profile',
      },
      {
        name: 'response_type',
        value: 'code',
      },
    ],
  },
  vk: {
    url: 'https://oauth.vk.com/authorize',
    queryParams: [
      {
        name: 'client_id',
        value: '7781054',
      },
      {
        name: 'redirect_uri',
        value: 'http://localhost:4000/api/auth/oauthVk',
      },
      {
        name: 'scope',
        value: 'offline+photos',
      },
      {
        name: 'state',
        value:
          'bdcbe28874ab05962b50430b1466a8ebcbda45ba8c3c1beee600699478ad2a4d',
      },
      {
        name: 'response_type',
        value: 'code',
      },
    ],
  },
  google: {
    url: 'https://accounts.google.com/o/oauth2/v2/auth',
    queryParams: [
      {
        name: 'client_id',
        value:
          '96975efecfd0e5efee67c9ac4cc350ac9372ae559b2fb8a08feba6841a33fb53',
      },
      {
        name: 'redirect_uri',
        value: 'http://localhost:4000/user/auth/oauthGoogle',
      },
      {
        name: 'scope',
        value: 'https://www.googleapis.com/auth/drive.metadata.readonly',
      },
      {
        name: 'state',
        value:
          'bdcbe28874ab05962b50430b1466a8ebcbda45ba8c3c1beee600699478ad2a4d',
      },
      {
        name: 'response_type',
        value: 'token',
      },
    ],
  },
};

const Icon42: React.FC = () => (
  <img
    height="30"
    width="30"
    src="https://signin.intra.42.fr/assets/42_logo-7dfc9110a5319a308863b96bda33cea995046d1731cebb735e41b16255106c12.svg"
    alt="42 auth"
    style={{ filter: 'invert(1)' }}
  />
);

const IconFacebook: React.FC = () => (
  <svg
    height="30"
    width="30"
    viewBox="0 0 409 409"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path
      fill="#475993"
      d="M353.701,0H55.087C24.665,0,0.002,24.662,0.002,55.085v298.616c0,30.423,24.662,55.085,55.085,55.085
	h147.275l0.251-146.078h-37.951c-4.932,0-8.935-3.988-8.954-8.92l-0.182-47.087c-0.019-4.959,3.996-8.989,8.955-8.989h37.882
	v-45.498c0-52.8,32.247-81.55,79.348-81.55h38.65c4.945,0,8.955,4.009,8.955,8.955v39.704c0,4.944-4.007,8.952-8.95,8.955
	l-23.719,0.011c-25.615,0-30.575,12.172-30.575,30.035v39.389h56.285c5.363,0,9.524,4.683,8.892,10.009l-5.581,47.087
	c-0.534,4.506-4.355,7.901-8.892,7.901h-50.453l-0.251,146.078h87.631c30.422,0,55.084-24.662,55.084-55.084V55.085
	C408.786,24.662,384.124,0,353.701,0z"
    />
  </svg>
);

const IconVk: React.FC = () => (
  <svg
    height="30"
    width="30"
    viewBox="0 0 112 112"
    xmlns="http://www.w3.org/2000/svg"
  >
    <g>
      <circle
        id="XMLID_11_"
        fill="#4D76A1"
        cx="56.098"
        cy="56.098"
        r="56.098"
      />
      <path
        fill="#FFFFFF"
        fillRule="evenodd"
        clipRule="evenodd"
        d="M53.979,80.702h4.403c0,0,1.33-0.146,2.009-0.878
		c0.625-0.672,0.605-1.934,0.605-1.934s-0.086-5.908,2.656-6.778c2.703-0.857,6.174,5.71,9.853,8.235
		c2.782,1.911,4.896,1.492,4.896,1.492l9.837-0.137c0,0,5.146-0.317,2.706-4.363c-0.2-0.331-1.421-2.993-7.314-8.463
		c-6.168-5.725-5.342-4.799,2.088-14.702c4.525-6.031,6.334-9.713,5.769-11.29c-0.539-1.502-3.867-1.105-3.867-1.105l-11.076,0.069
		c0,0-0.821-0.112-1.43,0.252c-0.595,0.357-0.978,1.189-0.978,1.189s-1.753,4.667-4.091,8.636c-4.932,8.375-6.904,8.817-7.71,8.297
		c-1.875-1.212-1.407-4.869-1.407-7.467c0-8.116,1.231-11.5-2.397-12.376c-1.204-0.291-2.09-0.483-5.169-0.514
		c-3.952-0.041-7.297,0.012-9.191,0.94c-1.26,0.617-2.232,1.992-1.64,2.071c0.732,0.098,2.39,0.447,3.269,1.644
		c1.135,1.544,1.095,5.012,1.095,5.012s0.652,9.554-1.523,10.741c-1.493,0.814-3.541-0.848-7.938-8.446
		c-2.253-3.892-3.954-8.194-3.954-8.194s-0.328-0.804-0.913-1.234c-0.71-0.521-1.702-0.687-1.702-0.687l-10.525,0.069
		c0,0-1.58,0.044-2.16,0.731c-0.516,0.611-0.041,1.875-0.041,1.875s8.24,19.278,17.57,28.993
		C44.264,81.287,53.979,80.702,53.979,80.702L53.979,80.702z"
      />
    </g>
  </svg>
);

const IconGoogle: React.FC = () => (
  <svg
    height="30"
    width="30"
    viewBox="0 0 512 512"
    xmlns="http://www.w3.org/2000/svg"
  >
    <g fill="none" fillRule="evenodd">
      <path
        d="M482.56 261.36c0-16.73-1.5-32.83-4.29-48.27H256v91.29h127.01c-5.47 29.5-22.1 54.49-47.09 71.23v59.21h76.27c44.63-41.09 70.37-101.59 70.37-173.46z"
        fill="#4285f4"
      ></path>
      <path
        d="M256 492c63.72 0 117.14-21.13 156.19-57.18l-76.27-59.21c-21.13 14.16-48.17 22.53-79.92 22.53-61.47 0-113.49-41.51-132.05-97.3H45.1v61.15c38.83 77.13 118.64 130.01 210.9 130.01z"
        fill="#34a853"
      ></path>
      <path
        d="M123.95 300.84c-4.72-14.16-7.4-29.29-7.4-44.84s2.68-30.68 7.4-44.84V150.01H45.1C29.12 181.87 20 217.92 20 256c0 38.08 9.12 74.13 25.1 105.99l78.85-61.15z"
        fill="#fbbc05"
      ></path>
      <path
        d="M256 113.86c34.65 0 65.76 11.91 90.22 35.29l67.69-67.69C373.03 43.39 319.61 20 256 20c-92.25 0-172.07 52.89-210.9 130.01l78.85 61.15c18.56-55.78 70.59-97.3 132.05-97.3z"
        fill="#ea4335"
      ></path>
      <path d="M20 20h472v472H20V20z"></path>
    </g>
  </svg>
);

const OAuthButton: React.FC<Props> = ({ variant }) => {
  const classes = useStyles();
  const [loading, setLoading] = useState(false);

  const handleClick = (e: React.MouseEvent<HTMLButtonElement>) => {
    setLoading(true);
    const queryParams = oauth[variant].queryParams.reduce((acc: any, el) => {
      acc[el.name] = el.value;
      return acc;
    }, {});
    const queryString = new URLSearchParams(queryParams).toString();
    const url = `${oauth[variant].url}?${queryString}`;
    window.location.href = url;
  };

  const icon = useMemo(() => {
    switch (variant) {
      case '42':
        return <Icon42 />;
      case 'facebook':
        return <IconFacebook />;
      case 'vk':
        return <IconVk />;
      case 'google':
        return <IconGoogle />;
    }
  }, [variant]);

  return (
    <Button
      type="submit"
      disabled={loading}
      onClick={handleClick}
      className={classes.OAuthButton}
      aria-label={`${variant} oauth`}
    >
      {loading ? <CircularProgress /> : icon}
    </Button>
  );
};

export default OAuthButton;
