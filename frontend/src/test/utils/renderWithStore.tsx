import { configureStore } from '@reduxjs/toolkit';
import { Provider } from 'react-redux';
import { render } from '@testing-library/react';
import authReducer from '@/features/auth/authSlice';

/**
 * Render a component with a redux store.
 *
 * By default, it will set up a store with the `authReducer` and no preloaded state.
 * You can override this by passing a `preloadedState` option, or a fully configured
 * `store` option.
 *
 * It will return an object with the store and the result of calling `render` with
 * the component wrapped in a `Provider`.
 */
export const renderWithStore = (ui: React.ReactElement, { preloadedState = {}, store = configureStore({
  reducer: {
    auth: authReducer,
  },
  preloadedState,
}) } = {}) => {
  return {
    store,
    ...render(<Provider store={store}>{ui}</Provider>),
  };
};
