import { useState, useCallback } from 'react';
import { produce } from 'immer';

/**
 * useState-like hook that uses immer for immutable updates
 * @param {*} initialState - Initial state value
 * @returns {Array} - [state, updateState] tuple
 *
 * @example
 * const [state, updateState] = useImmerState({ count: 0 });
 * updateState(draft => { draft.count += 1; });
 */
export function useImmerState(initialState) {
  const [state, setState] = useState(initialState);

  const updateState = useCallback((updater) => {
    setState((currentState) => produce(currentState, updater));
  }, []);

  return [state, updateState];
}
