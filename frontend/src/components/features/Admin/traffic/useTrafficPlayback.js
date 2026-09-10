import { useEffect, useState } from 'react';
import { clamp, DEFAULT_SLICE_MS, sanitizeSlice } from './trafficPresentation.js';

const PLAYBACK_INTERVAL_MS = 900;

export default function useTrafficPlayback({ isOpen, containerId, scopeKey, topology, isFetching }) {
  const [shift, setShift] = useState(0);
  const [slice, setSlice] = useState(DEFAULT_SLICE_MS);
  const [sliceInput, setSliceInput] = useState(String(DEFAULT_SLICE_MS));
  const [isPlaying, setIsPlaying] = useState(false);
  const windowInfo = topology?.window || { start: 0, end: 0, duration: slice, total: 0 };
  const totalDuration = Math.max(windowInfo.total || 0, topology?.total_duration || 0);
  const maxShift = Math.max(0, totalDuration - slice);
  const playbackFrames = Math.max(1, Math.floor(maxShift / Math.max(slice, 1)) + 1);
  const playbackIndex = Math.min(playbackFrames, Math.floor(Math.max(shift, 0) / Math.max(slice, 1)) + 1);

  useEffect(() => {
    setIsPlaying(false);
    if (!isOpen) return;
    setSlice(DEFAULT_SLICE_MS);
    setSliceInput(String(DEFAULT_SLICE_MS));
  }, [isOpen, scopeKey]);

  useEffect(() => {
    if (topology) setShift((current) => Math.min(current, maxShift));
  }, [maxShift, topology]);

  useEffect(() => {
    setSliceInput(String(slice));
  }, [slice]);

  useEffect(() => {
    if (!isPlaying || !isOpen || !containerId) return;
    if (maxShift <= 0 || shift >= maxShift) {
      setIsPlaying(false);
      return;
    }
    if (isFetching) return;
    const timer = window.setTimeout(() => {
      setShift((current) => Math.min(current + slice, maxShift));
    }, PLAYBACK_INTERVAL_MS);
    return () => window.clearTimeout(timer);
  }, [containerId, isFetching, isOpen, isPlaying, maxShift, shift, slice]);

  const stepPlayback = (direction) => {
    setIsPlaying(false);
    setShift((current) => clamp(current + direction * slice, 0, maxShift));
  };
  const handleSliceInputChange = (event) => setSliceInput(event.target.value);
  const commitSliceInput = () => {
    const nextSlice = sanitizeSlice(sliceInput);
    setSliceInput(String(nextSlice));
    setSlice(nextSlice);
    setIsPlaying(false);
  };
  const togglePlayback = () => {
    if (isPlaying) {
      setIsPlaying(false);
      return;
    }
    if (shift >= maxShift && maxShift > 0) setShift(0);
    setIsPlaying(true);
  };

  return {
    shift,
    setShift,
    slice,
    sliceInput,
    isPlaying,
    setIsPlaying,
    windowInfo,
    totalDuration,
    maxShift,
    playbackFrames,
    playbackIndex,
    canStepBackward: shift > 0,
    canStepForward: shift < maxShift,
    stepPlayback,
    handleSliceInputChange,
    commitSliceInput,
    togglePlayback,
  };
}
