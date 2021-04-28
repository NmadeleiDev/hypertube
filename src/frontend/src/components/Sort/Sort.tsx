import React from 'react';
import { SortRounded } from '@material-ui/icons';

import { useAppDispatch } from '../../store/store';
import { setSortingBy } from '../../store/features/UISlice';
import Dropdown, { IItem } from '../Dropdown/Dropdown';
import { useTranslation } from 'react-i18next';

const Sort = () => {
  const dispatch = useAppDispatch();
  const { t } = useTranslation();
  const onClick = (e: React.MouseEvent<HTMLButtonElement>) =>
    dispatch(setSortingBy({ sortBy: e.currentTarget.name }));

  const items = [
    {
      text: t('By rating'),
      name: 'rating',
      onClick,
    },
    {
      text: t('By name'),
      name: 'name',
      onClick,
    },
    {
      text: t('By year'),
      name: 'year',
      onClick,
    },
    {
      text: t('By avalibility'),
      name: 'avalibility',
      onClick,
    },
  ] as IItem[];

  return <Dropdown heroText={t('Sort')} icon={<SortRounded />} items={items} />;
};

export default Sort;
