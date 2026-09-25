// Cover the backend's four-minute startup and thirty-minute attachment queue
// deadlines, with room for the final status response. Stop tasks allow two minutes.
export const POLL_TIMEOUT = {
  running: 5 * 60 * 1000,
  stopped: 3 * 60 * 1000,
  attachment: 31 * 60 * 1000,
};
