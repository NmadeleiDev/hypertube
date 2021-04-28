import React from 'react';
// also exported from '@storybook/react' if you can deal with breaking changes in 6.1
import { Story, Meta } from '@storybook/react/types-6-0';

import CardsSlider, { CardsSliderProps } from './CardsSlider';

export default {
  title: 'CardsSlider',
  component: CardsSlider,
} as Meta;

const Template: Story<CardsSliderProps> = (args) => <CardsSlider {...args} />;

export const Default = Template.bind({});
