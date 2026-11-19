import { useCallback, useState } from 'react';
import { produce } from 'immer';

/** useState-like hook that applies immer draft updates. */
export function useImmerState(initialState) {
  const [state, setState] = useState(initialState);

  const updateState = useCallback((updater) => {
    setState((currentState) => produce(currentState, updater));
  }, []);

  return [state, updateState];
}
