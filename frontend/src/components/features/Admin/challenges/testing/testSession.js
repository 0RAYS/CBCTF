export const normalizeInstanceStatus = (status) => {
  const normalized = typeof status === 'string' ? status.toLowerCase() : '';
  return ['waiting', 'pending', 'terminating', 'running'].includes(normalized) ? normalized : '';
};

export const isInstanceTransitioning = (status) => ['waiting', 'pending', 'terminating'].includes(status);

export const hasReachedTarget = (status, target) =>
  target === 'running' ? status === 'running' : status !== 'running' && !isInstanceTransitioning(status);

export const initialTestState = () => ({
  testStatus: null,
  receivedAt: 0,
  loading: { status: false, starting: false, stopping: false },
});

// Each mounted dialog owns one session. Invalidating it also rejects in-flight responses.
export function createTestSession({ challengeId, api, onChange, notify }) {
  let state = initialTestState();
  let disposed = false;
  let generation = 0;
  let pollTimer;
  let deadlineTimer;

  const isCurrent = (version) => !disposed && version === generation;
  const publish = (patch) => {
    state = { ...state, ...patch };
    onChange(state);
  };
  const setLoading = (patch) => publish({ loading: { ...state.loading, ...patch } });
  const receive = (data) => publish({ testStatus: data, receivedAt: Date.now() });
  const invalidate = () => {
    clearTimeout(pollTimer);
    clearTimeout(deadlineTimer);
    return ++generation;
  };
  const finishPolling = () => {
    invalidate();
    setLoading({ status: false, starting: false, stopping: false });
  };

  function startPolling(target, version, immediate = false) {
    // The deadline is independent of request completion, including hung requests.
    deadlineTimer = setTimeout(
      () => {
        if (isCurrent(version)) finishPolling();
      },
      3 * 60 * 1000
    );

    async function poll(showLoading = false) {
      if (!isCurrent(version)) return;
      if (showLoading) setLoading({ status: true });
      try {
        const response = await api.getTestChallengeStatus(challengeId);
        if (!isCurrent(version)) return;
        if (response.code === 200) {
          receive(response.data);
          if (hasReachedTarget(normalizeInstanceStatus(response.data.remote?.status), target)) {
            finishPolling();
            return;
          }
        }
      } catch (error) {
        // Transient polling failures remain silent and retry until the deadline.
        if (showLoading && isCurrent(version)) notify('danger', 'fetchStatusFailed', error);
      } finally {
        if (showLoading && isCurrent(version)) setLoading({ status: false });
      }
      if (isCurrent(version)) pollTimer = setTimeout(poll, 5000);
    }

    if (immediate) return poll(true);
    pollTimer = setTimeout(poll, 5000);
  }

  async function load() {
    if (disposed || !challengeId) return;
    const version = invalidate();
    setLoading({ status: true, starting: false, stopping: false });
    try {
      const response = await api.getTestChallengeStatus(challengeId);
      if (!isCurrent(version)) return;
      if (response.code === 200) {
        receive(response.data);
        const status = normalizeInstanceStatus(response.data.remote?.status);
        if (isInstanceTransitioning(status)) {
          startPolling(status === 'terminating' ? 'stopped' : 'running', version);
        }
      }
    } catch (error) {
      if (isCurrent(version)) notify('danger', 'fetchStatusFailed', error);
    } finally {
      if (isCurrent(version)) setLoading({ status: false });
    }
  }

  async function act(action) {
    if (disposed || !challengeId || Object.values(state.loading).some(Boolean)) return;
    const version = invalidate();
    const starting = action === 'start';
    setLoading({ starting, stopping: !starting });
    try {
      const response = await (starting ? api.startTestVictim(challengeId) : api.stopTestVictim(challengeId));
      if (!isCurrent(version)) return;
      if (response.code !== 200) {
        setLoading({ starting: false, stopping: false });
        return;
      }
      if (starting) {
        receive({
          ...state.testStatus,
          remote: {
            ...state.testStatus?.remote,
            status: 'waiting',
            target: state.testStatus?.remote?.target || [],
            duration: state.testStatus?.remote?.duration || 0,
            remaining: 0,
          },
        });
        notify('success', 'actionSuccess');
        startPolling('running', version);
      } else {
        // Stop is synchronous on the backend; refresh immediately, then confirm if needed.
        notify('success', 'actionSuccess');
        await startPolling('stopped', version, true);
      }
    } catch (error) {
      if (!isCurrent(version)) return;
      notify('danger', 'actionFailed', error);
      setLoading({ starting: false, stopping: false });
    }
  }

  return {
    load,
    start: () => act('start'),
    stop: () => act('stop'),
    dispose: () => {
      disposed = true;
      invalidate();
    },
  };
}
