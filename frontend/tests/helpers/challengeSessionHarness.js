export { mountHook, useEffect, useRef, useState, useTranslation } from './pollingHookHarness.js';
export { toast } from '../../src/utils/toast.js';

export const api = {};
export const getChallengeStatus = (...args) => api.status(...args);
export const initChallenge = (...args) => api.init(...args);
export const resetChallenge = (...args) => api.reset(...args);
export const startRemoteTarget = (...args) => api.start(...args);
export const stopContainer = (...args) => api.stop(...args);
export const extendContainerTime = (...args) => api.extend(...args);
export const submitFlag = (...args) => api.submit(...args);
export const downloadChallengeAttachment = (...args) => api.download(...args);
export const downloadBlobResponse = (...args) => api.save(...args);
