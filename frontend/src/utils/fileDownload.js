/**
 * Download a blob response as a file.
 * Parses content-disposition header for filename if available.
 *
 * @param {Object} response - Axios response object (with response.data as blob/arraybuffer)
 * @param {string} [fallbackFilename='download'] - Fallback filename if not in headers
 * @param {string} [mimeType] - Optional MIME type for the blob (e.g., 'application/octet-stream')
 */
export function downloadBlobResponse(response, fallbackFilename = 'download', mimeType) {
  const contentDisposition = response.headers?.['content-disposition'];
  let filename = fallbackFilename;

  if (contentDisposition) {
    const matches = /filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/.exec(contentDisposition);
    if (matches != null && matches[1]) {
      filename = matches[1].replace(/['"]/g, '');
    }
  }

  const blobOptions = mimeType ? { type: mimeType } : undefined;
  const blob = response.data instanceof Blob ? response.data : new Blob([response.data], blobOptions);
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  link.remove();
  window.URL.revokeObjectURL(url);
}
