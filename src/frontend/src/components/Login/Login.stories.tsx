import React from 'react';
import { Story, Meta } from '@storybook/react/types-6-0';

import Login from './Login';

export default {
  title: 'Login',
  component: Login,
} as Meta;

const Template: Story = (args) => <Login {...args} />;

export const Default = Template.bind({});

export const OnColor = () => (
  <div
    style={{
      width: '100vw',
      height: '100vh',
      textAlign: 'center',
      paddingTop: '30vh',
      backdropFilter: 'blur(40px)',
      border: 'solid 2px transparent',
      backgroundClip: 'padding-box',
      background:
        'url(https://cdn.ananasposter.ru/image/cache/catalog/poster/film/99/1444-1000x830.jpg)',
    }}
  >
    <Login />
  </div>
);
