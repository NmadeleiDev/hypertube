import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface SortingItem {
  sortBy: string | null;
}

export interface UIState {
  sortBy: string | null;
  view: 'grid' | 'lines';
}

const initialState = {
  sortBy: null,
  view: 'lines',
} as UIState;

const UISlice = createSlice({
  name: 'UI',
  initialState,
  reducers: {
    setSortingBy(state, { payload }: PayloadAction<SortingItem>) {
      state.sortBy = payload.sortBy;
    },
    changeView(state) {
      state.view = state.view === 'grid' ? 'lines' : 'grid';
    },
  },
});

export const { setSortingBy, changeView } = UISlice.actions;

export default UISlice.reducer;
