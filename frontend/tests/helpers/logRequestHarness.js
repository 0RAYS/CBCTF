export const requests = [];

// The real API wrappers receive request.js's resolved business envelope, not an Axios response.
// Global request-interceptor notifications are outside these hook-session tests.
export default function request(config) {
  let resolve;
  let reject;
  const promise = new Promise((yes, no) => {
    resolve = yes;
    reject = no;
  });
  requests.push({ config, resolve, reject });
  return promise;
}
